package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type WriteReader interface {
	Write(*EnrichedSong, string) (int, error)
}

type Enricher interface {
	Update(*Song) (*http.Response, error)
}

type Service struct {
	writeData WriteReader
	findData  Enricher
}

func New(w WriteReader, e Enricher) *Service {
	return &Service{
		writeData: w,
		findData:  e,
	}
}

func (s *Service) ProseccAddSongRequest(r *http.Request, requestID string) ([]byte, int) {
	defer closeRequestBody(r.Body)

	req, statusCode, errResponse := validateAddSongRequest(r, requestID)
	if statusCode != http.StatusOK {
		return errResponse, statusCode
	}

	response, statusCode := s.processAddSong(req, requestID)
	return response, statusCode
}

func closeRequestBody(b io.ReadCloser) {
	if b != nil {
		_ = b.Close()
	}
}

func validateAddSongRequest(r *http.Request, requestID string) (*Song, int, []byte) {
	if r.Method != http.MethodPost {
		msg := fmt.Sprintf(errMethod, http.MethodPost, r.Method)
		dataJson, statusCode := createAddSongResponse(
			false, http.StatusMethodNotAllowed, msg, requestID, nil)
		return nil, statusCode, dataJson
	}

	req, err := decodeBodyAddSongRequest(r.Body)
	if err != nil {
		dataJson, statusCode := createAddSongResponse(
			false, http.StatusBadRequest, fmt.Sprint(err), requestID, nil)
		return nil, statusCode, dataJson
	}

	if err = validateBodyAddSongRequest(req); err != nil {
		dataJson, statusCode := createAddSongResponse(
			false, http.StatusBadRequest, fmt.Sprint(err), requestID, nil)
		return nil, statusCode, dataJson
	}
	return req, http.StatusOK, nil
}

func createAddSongResponse(ok bool, statusCode int, msg, requestID string, id *int) ([]byte, int) {
	log.Printf("[%s]  %s\n", requestID, msg)
	resp := ResponsePost{
		Sucsess:    ok,
		Message:    msg,
		StatusCode: statusCode,
	}
	if id != nil {
		resp.ID = id
	}
	dataJson, err := json.Marshal(resp)
	if err != nil {
		log.Println(errMarshalJson)
		return nil, http.StatusInternalServerError
	}
	return dataJson, statusCode
}

func decodeBodyAddSongRequest(b io.ReadCloser) (*Song, error) {
	var s Song
	err := json.NewDecoder(b).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func validateBodyAddSongRequest(req *Song) error {
	const (
		song         = "song"
		group        = "group"
		maxSymbSong  = 200
		maxSymbGroup = 100
	)
	var bufErrors bytes.Buffer

	if req.Group == "" {
		bufErrors.WriteString(fmt.Sprintf("|пустое поле `%s`|", group))
	}
	if req.Song == "" {
		bufErrors.WriteString(fmt.Sprintf("|пустое поле `%s`|", song))
	}

	symbGroup := []byte(req.Group)
	if len(symbGroup) > maxSymbGroup {
		bufErrors.WriteString(
			fmt.Sprintf("|для поля `%s` много символов: %d|", group, symbGroup))
	}
	symbSong := []byte(req.Song)
	if len(symbSong) > maxSymbSong {
		bufErrors.WriteString(
			fmt.Sprintf("|для поля `%s` много символов: %d|", song, symbSong))
	}

	var err error
	if bufErrors.Len() > 0 {
		err = errors.New(bufErrors.String())
	}

	return err
}
