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
	"time"
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

func (s *Service) CheckAddSong(r *http.Request, requestID string) ([]byte, int) {
	defer func() {
		if r.Body != nil {
			_ = r.Body.Close()
		}
	}()

	if r.Method != http.MethodPost {
		msg := fmt.Sprintf(errMethod, http.MethodPost, r.Method)
		dataJson, statusCode := prepareResponseAddSong(
			false, http.StatusMethodNotAllowed, msg, requestID, nil)
		return dataJson, statusCode
	}

	req, err := parseRequestServer(r.Body)
	if err != nil {
		dataJson, statusCode := prepareResponseAddSong(
			false, http.StatusBadRequest, fmt.Sprint(err), requestID, nil)
		return dataJson, statusCode
	}

	err = validateBodyRequest(req)
	if err != nil {
		dataJson, statusCode := prepareResponseAddSong(
			false, http.StatusBadRequest, fmt.Sprint(err), requestID, nil)
		return dataJson, statusCode
	}

	dataRespClient, err := s.findData.Update(req)
	if err != nil {
		dataJson, statusCode := prepareResponseAddSong(
			false, http.StatusInternalServerError, fmt.Sprint(err), requestID, nil)
		return dataJson, statusCode
	}
	dataClient, err := parseRequestClient(dataRespClient.Body)
	if err != nil {
		dataJson, statusCode := prepareResponseAddSong(
			false, http.StatusInternalServerError, fmt.Sprint(err), requestID, nil)
		return dataJson, statusCode
	}
	err = validateBodyResponseClient(dataClient)
	if err != nil {
		dataJson, statusCode := prepareResponseAddSong(
			false, http.StatusInternalServerError, fmt.Sprint(err), requestID, nil)
		return dataJson, statusCode
	}
	enrichedData := enrichedData(req, dataClient)
	idSong, err := s.writeData.Write(enrichedData, requestID)
	if err != nil {
		dataJson, statusCode := prepareResponseAddSong(
			false, http.StatusInternalServerError, fmt.Sprint(err), requestID, nil)
		return dataJson, statusCode
	}

	dataJson, statusCode := prepareResponseAddSong(
		true, http.StatusInternalServerError, msg201, requestID, &idSong)
	return dataJson, statusCode
}

func prepareResponseAddSong(ok bool, statusCode int, msg, requestID string, id *int) ([]byte, int) {
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

func parseRequestServer(b io.ReadCloser) (*Song, error) {
	var s Song
	err := json.NewDecoder(b).Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func validateBodyRequest(req *Song) error {
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

func parseRequestClient(b io.ReadCloser) (*EnrichedSong, error) {
	var e EnrichedSong
	err := json.NewDecoder(b).Decode(&e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func validateBodyResponseClient(resp *EnrichedSong) error {
	const (
		layout = "02.01.2006"
	)

	var bufErrors bytes.Buffer

	parsedDate, err := time.Parse(layout, resp.ReleaseDate)
	if err != nil {
		bufErrors.WriteString(fmt.Sprintf("|дата не соответствует шаблону. Имеется: `%s`, должно быть: `%s`|",
			resp.ReleaseDate, layout))
	}
	now := time.Now()
	if err == nil && parsedDate.After(now) {
		bufErrors.WriteString(fmt.Sprintf("|дата песни не должна быть в будущем. Имеется: `%s`, сейчас: `%s`|",
			resp.ReleaseDate, now.Format(layout)))
	}

	if _, err = url.Parse(resp.Link); err != nil {
		bufErrors.WriteString(fmt.Sprintf("|нет ссылки на песню. Имеется: `%s`|",
			resp.Link))
	}

	if bufErrors.Len() > 0 {
		err = errors.New(bufErrors.String())
	}

	return err
}

func enrichedData(s *Song, e *EnrichedSong) *EnrichedSong {
	const (
		layout = "02.01.2006"
	)
	parsedDate, _ := time.Parse(layout, e.ReleaseDate)
	e.Group = s.Group
	e.Song = s.Song
	e.ReleaseDateTime = parsedDate
	return e
}
