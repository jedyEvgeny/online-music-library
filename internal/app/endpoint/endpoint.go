package endpoint

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

type Checker interface {
	CheckAddSong(*http.Request, string) ([]byte, int)
	// CheckGet(*http.Request, string) ([]byte, int)
}

type Endpoint struct {
	checkRequest Checker
}

func New(c Checker) *Endpoint {
	return &Endpoint{
		checkRequest: c,
	}
}

const (
	msgRequest = "[%s] Получен запрос с методом: %s от URL: %s\n"
)

func (e *Endpoint) HandlerAddSong(w http.ResponseWriter, r *http.Request) {
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, status := e.checkRequest.CheckAddSong(r, reqID)
	w.WriteHeader(status)
	w.Write(resp)
}

func requestID() string {
	return uuid.New().String()
}
