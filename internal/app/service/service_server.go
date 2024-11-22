package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type WriteReader interface {
	Write(*EnrichedSong, string) (int, error)
}

type DelUpdater interface {
	Delete(int, string) error
}

type Enricher interface {
	Update(*Song) (*http.Response, error)
}

type Service struct {
	writeData     WriteReader
	findData      Enricher
	delUpdateData DelUpdater
}

func New(w WriteReader, e Enricher, d DelUpdater) *Service {
	return &Service{
		writeData:     w,
		findData:      e,
		delUpdateData: d,
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

func (s *Service) ProseccDelSongRequest(r *http.Request, requestID string) ([]byte, int) {
	defer closeRequestBody(r.Body)

	id, statusCode, errResponse := validateDelSongRequest(r, requestID)
	if statusCode != http.StatusNoContent {
		return errResponse, statusCode
	}

	response, statusCode := s.delSongFromStorage(id, requestID)
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

func validateDelSongRequest(r *http.Request, requestID string) (int, int, []byte) {
	if r.Method != http.MethodDelete {
		msg := fmt.Sprintf(errMethod, http.MethodDelete, r.Method)
		dataJson, statusCode := createDelSongResponse(
			false, http.StatusMethodNotAllowed, msg, requestID)
		return 0, statusCode, dataJson
	}

	idStr := decodeDelSongRequest(r.URL)

	id, err := validateParamDelSongRequest(idStr)
	if err != nil {
		dataJson, statusCode := createDelSongResponse(
			false, http.StatusBadRequest, fmt.Sprint(err), requestID)
		return 0, statusCode, dataJson
	}
	return id, http.StatusNoContent, nil
}

func createDelSongResponse(ok bool, statusCode int, msg, requestID string) ([]byte, int) {
	log.Printf("[%s]  %s\n", requestID, msg)
	resp := ResponseDelete{
		Sucsess:    ok,
		Message:    msg,
		StatusCode: statusCode,
	}

	dataJson, err := json.Marshal(resp)
	if err != nil {
		log.Printf(errMarshalJson, err)
		return nil, http.StatusInternalServerError
	}
	return dataJson, statusCode
}

func decodeDelSongRequest(u *url.URL) string {
	query := u.Query()
	return query.Get("song_id")
}

func validateParamDelSongRequest(idStr string) (int, error) {
	songID, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf(errIDDel, err)
	}

	if songID <= 0 {
		return 0, fmt.Errorf(errIDValueDel, err)
	}
	return songID, nil
}

func (s *Service) delSongFromStorage(id int, requestID string) ([]byte, int) {
	err := s.delUpdateData.Delete(id, requestID)

	if err != nil {
		dataJson, statusCode := createDelSongResponse(
			false, http.StatusInternalServerError, fmt.Sprint(err), requestID)
		return dataJson, statusCode
	}

	return nil, http.StatusNoContent
}
