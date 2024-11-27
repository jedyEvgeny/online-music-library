package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func (s *Service) processAddSong(req *Song, requestID string) ([]byte, int) {
	dataClient, statusCode, errResponse := s.enrighedSong(req, requestID)
	if statusCode != http.StatusOK {
		return errResponse, statusCode
	}
	dataJson, statusCode := s.saveSongToStorage(req, dataClient, requestID)

	return dataJson, statusCode
}

func (s *Service) enrighedSong(req *Song, requestID string) (*EnrichedSong, int, []byte) {
	dataRespClient, err := s.enricher.Update(req, requestID)
	if err != nil {
		dataJson, statusCode := s.createAddSongResponse(
			false, http.StatusInternalServerError, fmt.Sprint(err), requestID, nil)
		return nil, statusCode, dataJson
	}

	dataClient, err := decodeBodyResponse(dataRespClient.Body)
	if err != nil {
		dataJson, statusCode := s.createAddSongResponse(
			false, http.StatusInternalServerError, fmt.Sprint(err), requestID, nil)
		return nil, statusCode, dataJson
	}

	if err = validateEnrichedSongFields(dataClient); err != nil {
		dataJson, statusCode := s.createAddSongResponse(
			false, http.StatusInternalServerError, fmt.Sprint(err), requestID, nil)
		return nil, statusCode, dataJson
	}
	return dataClient, http.StatusOK, nil
}

func decodeBodyResponse(b io.ReadCloser) (*EnrichedSong, error) {
	var e EnrichedSong
	err := json.NewDecoder(b).Decode(&e)
	if err != nil {
		return nil, err
	}
	return &e, nil
}

func validateEnrichedSongFields(resp *EnrichedSong) error {
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

func enrichSongData(s *Song, e *EnrichedSong) *EnrichedSong {
	e.Group = s.Group
	e.Song = s.Song
	return e
}

func (s *Service) saveSongToStorage(req *Song, dataClient *EnrichedSong, requestID string) ([]byte, int) {
	enrichedData := enrichSongData(req, dataClient)

	idSong, err := s.repository.Write(enrichedData, requestID)
	if err != nil {
		dataJson, statusCode := s.createAddSongResponse(false, http.StatusInternalServerError, fmt.Sprint(err), requestID, nil)
		return dataJson, statusCode
	}
	dataJson, statusCode := s.createAddSongResponse(
		true, http.StatusCreated, msg201, requestID, &idSong)
	return dataJson, statusCode
}
