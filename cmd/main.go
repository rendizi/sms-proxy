package main

import (
	"fmt"
	"github.com/rendizi/sms-proxy/internal/handler"
	"net/http"
	"os"

	"github.com/MadAppGang/httplog"
)

var (
	loginHandler      http.Handler = http.HandlerFunc(handler.Login)
	yearsHandler      http.Handler = http.HandlerFunc(handler.Years)
	termsHandler      http.Handler = http.HandlerFunc(handler.Terms)
	diaryHandler      http.Handler = http.HandlerFunc(handler.Diary)
	evaluationHandler http.Handler = http.HandlerFunc(handler.Evaluation)
)

func main() {
	mux := http.NewServeMux()

	loggerWithFormatter := httplog.LoggerWithFormatter(httplog.DefaultLogFormatterWithRequestHeader)
	mux.Handle("/login", loggerWithFormatter(loginHandler))
	mux.Handle("/years", loggerWithFormatter(yearsHandler))
	mux.Handle("/years/terms", loggerWithFormatter(termsHandler))
	mux.Handle("/diary", loggerWithFormatter(diaryHandler))
	mux.Handle("/diary/subjects/evaluation", loggerWithFormatter(evaluationHandler))

	corsHandler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	//слушаем на порту 8080
	err := http.ListenAndServe(":8080", corsHandler(mux))
	if err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("server closed")
		} else {
			fmt.Printf("error starting server: %s\n", err)
			os.Exit(1)
		}
	}
}
