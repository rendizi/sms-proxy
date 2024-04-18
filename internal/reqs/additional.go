package reqs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

type respopnse struct {
	Success bool `json:"success"`
	Data    []struct {
		Name string `json:"name"`
		Id   string `json:"id"`
	} `json:"data"`
	Message string `json:"message"`
}

func SetCookies(req *http.Request, appToken, sessionToken string) {
	req.Header.Set("Host", "sms.akb.nis.edu.kz")
	req.Header.Set("Accept", "/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://sms.akb.nis.edu.kz")
	req.Header.Set("Referer", "https://sms.akb.nis.edu.kz/jcediary/index/0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1.2 Mobile/15E148 Safari/604.1")

	cookie := fmt.Sprintf("UserSessionKey=%s; lang=ru-RU; ApplicationAuth=%s; Aktobe_SessionID=toinkccizpmpmkwi5so4uffg", sessionToken, appToken)
	req.Header.Set("Cookie", cookie)
}

func SetJceCookies(req *http.Request, appToken, sessionToken, url string) {
	req.Header.Set("Host", "sms.akb.nis.edu.kz")
	req.Header.Set("Accept", "/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://sms.akb.nis.edu.kz")
	req.Header.Set("Referer", url)
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1.2 Mobile/15E148 Safari/604.1")

	cookie := fmt.Sprintf("UserSessionKey=%s; lang=ru-RU; ApplicationAuth=%s; Aktobe_SessionID=toinkccizpmpmkwi5so4uffg; Aktobe_jce_SessionID=lacdu5m101ybdc3rz2xkirxm", sessionToken, appToken)
	req.Header.Set("Cookie", cookie)
}

func Parallels(appAuth, sessionToken, periodId string) (string, error) {
	smsUrl := "https://sms.akb.nis.edu.kz/JceDiary/GetParallels?_dc=1713329454662"

	data := url.Values{}
	data.Set("periodId", periodId)
	data.Set("page", "1")
	data.Set("start", "0")
	data.Set("limit", "100")

	paylod := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", smsUrl, paylod)
	if err != nil {
		return "", err
	}

	SetCookies(req, appAuth, sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respData respopnse

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return "", err
	}

	log.Println(respData)

	if !respData.Success {
		return "", errors.New(respData.Message)
	}

	return respData.Data[0].Id, nil
}

func Klasses(appAuth, sessionToken, periodId, parallelId string) (string, error) {
	smsUrl := "https://sms.akb.nis.edu.kz/JceDiary/GetKlasses"

	// Define form data
	formData := url.Values{}

	formData.Set("periodId", periodId)
	formData.Set("parallelId", parallelId)
	formData.Set("page", "1")
	formData.Set("start", "0")
	formData.Set("limit", "100")

	payload := bytes.NewBufferString(formData.Encode())

	// Create HTTP client
	client := &http.Client{}

	// Create POST request
	req, err := http.NewRequest("POST", smsUrl, payload)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "", err
	}

	// Set form data
	req.Form = formData

	SetCookies(req, appAuth, sessionToken)

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return "", err
	}
	defer resp.Body.Close()

	var respData respopnse

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		log.Println("Error decoding response:", err)
		return "", err
	}

	log.Println(respData)

	if respData.Message != "" {
		log.Println("Error message in response:", respData.Message)
		return "", errors.New(respData.Message + ". Status code:" + resp.Status)
	}

	klassID := respData.Data[0].Id
	log.Println("Klass ID:", klassID)

	return klassID, nil
}

func Students(appAuth, sessionToken, periodId, klassId string) (string, error) {
	smsUrl := "https://sms.akb.nis.edu.kz/JceDiary/GetStudents?_dc=1713329458147"

	data := url.Values{}
	data.Set("periodId", periodId)
	data.Set("klassId", klassId)
	data.Set("page", "1")
	data.Set("start", "0")
	data.Set("limit", "100")

	paylod := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", smsUrl, paylod)
	if err != nil {
		return "", err
	}

	SetCookies(req, appAuth, sessionToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var respData respopnse

	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return "", err
	}

	log.Println(respData)

	if !respData.Success {
		return "", errors.New(respData.Message)
	}

	return respData.Data[0].Id, nil
}
