package logger

import (
	"log"
	"net/http"
	"os"
	"time"
)

type Logger struct {
	log.Logger
}

func NewLogger() *Logger {
	l := new(Logger)

	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	l.SetOutput(file)

	return l
}

func (l *Logger) LogRequest(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l.Println("\n[" + time.Now().Format("2006-01-02 15:04:05") + "]")
		l.Printf("-----Request-----\nmethod: %s\nURI: %s\n", r.Method, r.RequestURI)
		handler.ServeHTTP(w, r)
	}
}
