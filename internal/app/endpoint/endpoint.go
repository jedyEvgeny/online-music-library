package endpoint

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

type delAdder interface {
	ProseccAddSongRequest(*http.Request, string) ([]byte, int)
	ProseccDelSongRequest(*http.Request, string) ([]byte, int)
}

type UpdateReader interface {
	ProcessUpdateSongRequest(*http.Request, string) ([]byte, int)
	ProcessReadLirycsSongRequest(*http.Request, string) ([]byte, int)
	ProcessLibraryRequest(*http.Request, string) ([]byte, int)
}

type Endpoint struct {
	process delAdder
	update  UpdateReader
}

func New(c delAdder, u UpdateReader) *Endpoint {
	return &Endpoint{
		process: c,
		update:  u,
	}
}

const (
	msgRequest = "[%s] Получен запрос с методом: %s от URL: %s\n"
)

func (e *Endpoint) HandlerAddSong(w http.ResponseWriter, r *http.Request) {
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, status := e.process.ProseccAddSongRequest(r, reqID)
	w.WriteHeader(status)
	w.Write(resp)
}

func (e *Endpoint) HandlerDeleteSong(w http.ResponseWriter, r *http.Request) {
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, status := e.process.ProseccDelSongRequest(r, reqID)
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

func (e *Endpoint) HandlerLiryc(w http.ResponseWriter, r *http.Request) {
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, status := e.update.ProcessReadLirycsSongRequest(r, reqID)
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

func (e *Endpoint) HandlerLibrary(w http.ResponseWriter, r *http.Request) {
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, status := e.update.ProcessLibraryRequest(r, reqID)
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

// Дописать
func (e *Endpoint) HandlerPatchSong(w http.ResponseWriter, r *http.Request) {
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, status := e.update.ProcessUpdateSongRequest(r, reqID)
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

func requestID() string {
	return uuid.New().String()
}
