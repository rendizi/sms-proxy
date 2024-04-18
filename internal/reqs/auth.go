package reqs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/url"
	_ "net/url"
	"regexp"
	"time"
)

//curl -X POST "https://sms.akb.nis.edu.kz/root/Account/LogOn" \
//     -H "Host: sms.akb.nis.edu.kz" \
//     -H "Accept: /" \
//     -H "X-Requested-With: XMLHttpRequest" \
//     -H "Sec-Fetch-Site: same-origin" \
//     -H "Accept-Language: en-GB,en-US;q=0.9,en;q=0.8" \
//     -H "Accept-Encoding: gzip, deflate, br" \
//     -H "Sec-Fetch-Mode: cors" \
//     -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8" \
//     -H "Origin: https://sms.akb.nis.edu.kz" \
//     -H "Referer: https://sms.akb.nis.edu.kz/Root/Account/Login?ReturnUrl=%2froot" \
//     -H "Connection: keep-alive" \
//     -H "Sec-Fetch-Dest: empty" \
//     -H "Cookie: lang=ru-RU; Aktobe_SessionID=toinkccizpmpmkwi5so4uffg" \
//     -H "User-Agent: Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1.2 Mobile/15E148 Safari/604.1" \
//     --data-urlencode "login=080930552551" \
//     --data-urlencode "password=Ali04xan!" \
//     --data-urlencode "twoFactorAuthCode=" \
//     --data-urlencode "captchaInput=" \
//     --data-urlencode "application2FACode="

var jwtSecret = []byte("secret_key")

func LogIn(login, password string) (string, string, error) {
	dest := "https://sms.akb.nis.edu.kz/root/Account/LogOn"

	data := url.Values{}
	data.Set("login", login)
	data.Set("password", password)
	data.Set("twoFactorAuthCode", "")
	data.Set("captchaInput", "")
	data.Set("application2FACode", "")

	payload := bytes.NewBufferString(data.Encode())

	req, err := http.NewRequest("POST", dest, payload)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Host", "sms.akb.nis.edu.kz")
	req.Header.Set("Accept", "/")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Origin", "https://sms.akb.nis.edu.kz")
	req.Header.Set("Referer", "https://sms.akb.nis.edu.kz/Root/Account/Login?ReturnUrl=%2froot")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Cookie", "lang=ru-RU; Aktobe_SessionID=toinkccizpmpmkwi5so4uffg")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_1_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1.2 Mobile/15E148 Safari/604.1")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var respData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return "", "", err
	}

	var userSessionKey, applicationAuth string
	for _, cookie := range resp.Header["Set-Cookie"] {
		if len(userSessionKey) == 0 {
			userSessionKey = extractCookieValue(cookie, "UserSessionKey")
		}
		if len(applicationAuth) == 0 {
			applicationAuth = extractCookieValue(cookie, "ApplicationAuth")

		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ApplicationAuth": applicationAuth,
		"UserSessionKey":  userSessionKey,
		"exp":             time.Now().Add(15 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString(jwtSecret)

	return respData["message"].(string), tokenString, nil
}

func extractCookieValue(setCookieHeader, cookieName string) string {
	regex := fmt.Sprintf(`%s=([^;]+)`, cookieName)
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(setCookieHeader)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}
