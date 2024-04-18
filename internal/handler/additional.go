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

type years struct {
	Success bool `json:"success"`
	Data    []struct {
		Name  string `json:"name"`
		DData struct {
			IsActual  bool   `json:"IsActual"`
			BeginDate string `json:"BeginDate"`
			EndDate   string `json:"EndDate"`
		} `json:"Data"`
		Id string `json:"id"`
	}
	Message string `json:"message"`
}

type terms struct {
	Success bool `json:"success"`
	Data    []struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"data"`
	Message string `json:"message"`
}

func Years(w http.ResponseWriter, r *http.Request) {
	appAuth, sessionsKey := server.GetTokens(w, r)
	if appAuth == "" || sessionsKey == "" {
		return
	}

	smsUrl := "https://sms.akb.nis.edu.kz/Ref/GetSchoolYears?fullData=true"

	data := url.Values{}
	data.Set("page", "1")
	data.Set("start", "0")
	data.Set("limit", "100")

	payload := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", smsUrl, payload)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	reqs.SetCookies(req, appAuth, sessionsKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}
	defer resp.Body.Close()

	log.Println(sessionsKey)

	var respData years

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

func Terms(w http.ResponseWriter, r *http.Request) {
	appAuth, sessionsKey := server.GetTokens(w, r)
	if appAuth == "" || sessionsKey == "" {
		return
	}

	yearId := r.URL.Query().Get("id")
	if yearId == "" {
		server.Error(map[string]interface{}{"message": "year id is not provided"}, w)
		return
	}

	smsUrl := "https://sms.akb.nis.edu.kz/Ref/GetPeriods"

	data := url.Values{}
	data.Set("schoolYearId", yearId)
	data.Set("page", "1")
	data.Set("start", "0")
	data.Set("limit", "100")

	payload := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", smsUrl, payload)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	reqs.SetCookies(req, appAuth, sessionsKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}
	defer resp.Body.Close()

	var respData terms

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	log.Println(respData)

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
