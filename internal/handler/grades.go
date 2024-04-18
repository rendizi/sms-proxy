package handler

import (
	"bytes"
	"encoding/json"
	"github.com/rendizi/sms-proxy/internal/reqs"
	"github.com/rendizi/sms-proxy/internal/server"
	"log"
	"net/http"
	"net/url"
)

type jceresp struct {
	Success bool `json:"success"`
	Data    struct {
		Url string `json:"url"`
	} `json:"data"`
	Message string `json:"message"`
}

type subjects struct {
	Success bool `json:"success"`
	Data    []struct {
		Name        string  `json:"name"`
		JournalId   string  `json:"JournalId"`
		Score       float32 `json:"score"`
		Mark        int     `json:"mark"`
		Evaluations []struct {
			Type      int                `json:"type"`
			Name      string             `json:"ShortName"`
			MaxScores map[string]float64 `json:"MaxScores"`
			Id        string             `json:"id"`
		} `json:"evaluations"`
		Id string `json:"id"`
	} `json:"data"`
	Message string `json:"message"`
	Url     string `json:"url"`
}

func Diary(w http.ResponseWriter, r *http.Request) {
	appAuth, sessionsKey := server.GetTokens(w, r)
	log.Println("AppAuth:", appAuth)
	log.Println("SessionsKey:", sessionsKey)
	if appAuth == "" || sessionsKey == "" {
		log.Println("AppAuth or SessionsKey is empty")
		return
	}
	periodId := r.URL.Query().Get("periodId")
	log.Println("PeriodId:", periodId)
	if periodId == "" {
		server.Error(map[string]interface{}{"message": "period id is not provided"}, w)
		return
	}

	parallelId, err := reqs.Parallels(appAuth, sessionsKey, periodId)
	if err != nil {
		log.Println("Error getting parallel:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}
	log.Println("ParallelId:", parallelId)

	klassId, err := reqs.Klasses(appAuth, sessionsKey, periodId, parallelId)
	if err != nil {
		log.Println("Error getting klass:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}
	log.Println("KlassId:", klassId)

	studentId, err := reqs.Students(appAuth, sessionsKey, periodId, klassId)
	if err != nil {
		log.Println("Error getting student:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}
	log.Println("StudentId:", studentId)

	smsUrl := "https://sms.akb.nis.edu.kz/JceDiary/GetJceDiary"
	log.Println("SMS URL:", smsUrl)

	data := url.Values{}
	data.Set("periodId", periodId)
	data.Set("parallelId", parallelId)
	data.Set("klassId", klassId)
	data.Set("studentId", studentId)
	log.Println("Data:", data)

	payload := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", smsUrl, payload)
	if err != nil {
		log.Println("Error creating request:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}
	log.Println("Request created")

	reqs.SetCookies(req, appAuth, sessionsKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error performing request:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}
	defer resp.Body.Close()

	var respData jceresp
	log.Println("Response received")

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		log.Println("Error decoding response:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	if !respData.Success {
		log.Println("Request unsuccessful:", respData.Message)
		server.Error(map[string]interface{}{"message": respData.Message}, w)
		return
	}
	//Zdes konec
	log.Println("-------Subjects-------")
	Subjects(w, appAuth, sessionsKey, respData.Data.Url)

}

func Subjects(w http.ResponseWriter, appAuth, sessionAuth, refereeUrl string) {
	smsUrl := "https://sms.akb.nis.edu.kz/Jce/Diary/GetSubjects?_dc=1713235818260"
	log.Println("New SMS URL:", smsUrl)

	data := url.Values{}
	data.Set("page", "1")
	data.Set("start", "0")
	data.Set("limit", "100")
	log.Println("New Data:", data)

	payload := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", smsUrl, payload)
	if err != nil {
		log.Println("Error creating new request:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}
	log.Println("New request created")

	reqs.SetJceCookies(req, appAuth, sessionAuth, refereeUrl)
	log.Println(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error performing new request:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}
	defer resp.Body.Close()

	var anotherRespData subjects
	log.Println("New response received")

	err = json.NewDecoder(resp.Body).Decode(&anotherRespData)
	if err != nil {
		log.Println("Error decoding new response:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	anotherRespData.Url = refereeUrl

	jsonData, err := json.Marshal(anotherRespData)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	log.Println("Sending response")
	log.Println(anotherRespData)

	_, err = w.Write(jsonData)
	if err != nil {
		log.Println("Error writing response:", err)
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type evaluationReq struct {
	JournalId string `json:"journalId"`
	EvalId    string `json:"evalId"`
	Url       string `json:"url"`
}

type evaluationResp struct {
	Success bool `json:"success"`
	Data    []struct {
		Name     string  `json:"name"`
		Score    float32 `json:"score"`
		MaxScore float32 `json:"maxscore"`
		RubricId string  `json:"RubricId"`
		Id       string  `json:"id"`
	}
	Message string `json:"message"`
}

func Evaluation(w http.ResponseWriter, r *http.Request) {
	var reqData evaluationReq
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	appAuth, sessionsKey := server.GetTokens(w, r)
	if appAuth == "" || sessionsKey == "" {
		return
	}
	if reqData.JournalId == "" || reqData.EvalId == "" || reqData.Url == "" {
		server.Error(map[string]interface{}{"message": "data is not provided"}, w)
		return
	}

	smsUrl := "https://sms.akb.nis.edu.kz/Ref/GetResultByEvalution"

	data := url.Values{}
	data.Set("journalId", reqData.JournalId)
	data.Set("evalId", reqData.EvalId)
	data.Set("page", "1")
	data.Set("start", "0")
	data.Set("limir", "100")

	payload := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", smsUrl, payload)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	reqs.SetJceCookies(req, appAuth, sessionsKey, reqData.Url)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}
	defer resp.Body.Close()

	var respData evaluationResp

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	jsonData, err := json.Marshal(respData)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonData)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}
