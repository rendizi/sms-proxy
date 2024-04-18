package handler

import (
	"encoding/json"
	"github.com/rendizi/sms-proxy/internal/reqs"
	"github.com/rendizi/sms-proxy/internal/server"
	"net/http"
)

type user struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds user
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		server.Error(map[string]interface{}{"message": err.Error()}, w)
		return
	}

	if len(creds.Login) != 12 || creds.Password == "" {
		server.Error(map[string]interface{}{"message": "invalid data"}, w)
		return
	}

	message, token, err := reqs.LogIn(creds.Login, creds.Password)
	if err != nil {
		server.Error(map[string]interface{}{"message": message}, w)
		return
	}

	server.Ok(map[string]interface{}{"message": message, "token": token}, w)
}
