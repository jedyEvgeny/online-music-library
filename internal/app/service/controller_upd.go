package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func (s *Service) ProcessUpdateSongRequest(r *http.Request, requestID string) ([]byte, int) {
	defer closeRequestBody(r.Body)

	data, idSong, statusCode, errResponse := validatePatchSongRequest(r, requestID)
	if statusCode != http.StatusOK {
		return errResponse, statusCode
	}
	response, statusCode := s.updSongInStorage(data, idSong, requestID)
	return response, statusCode
}

func validatePatchSongRequest(r *http.Request, requestID string) (*EnrichedSong, int, int, []byte) {
	if r.Method != http.MethodPatch {
		msg := fmt.Sprintf(errMethod, http.MethodPatch, r.Method)
		dataJson, statusCode := createPatchDelSongResponse(
			false, http.StatusMethodNotAllowed, msg, requestID)
		return nil, 0, statusCode, dataJson
	}

	req, idSong, err := decodeBodyPatchSongRequest(r)
	if err != nil {
		dataJson, statusCode := createPatchDelSongResponse(
			false, http.StatusBadRequest, fmt.Sprint(err), requestID)
		return nil, 0, statusCode, dataJson
	}

	if err := validateBodyPatchSongRequest(req); err != nil {
		dataJson, statusCode := createPatchDelSongResponse(
			false, http.StatusBadRequest, fmt.Sprint(err), requestID)
		return nil, 0, statusCode, dataJson
	}
	return req, idSong, http.StatusOK, nil
}

func decodeBodyPatchSongRequest(r *http.Request) (*EnrichedSong, int, error) {
	songIdStr := strings.TrimPrefix(r.URL.Path, "/song-upd/")
	songID, err := validateIdSong(songIdStr)
	if err != nil {
		return nil, 0, nil
	}

	e := EnrichedSong{}
	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		return nil, 0, err
	}
	return &e, songID, nil
}

func validateIdSong(songIdStr string) (int, error) {
	idSong, err := strconv.Atoi(songIdStr)
	if err != nil {
		return 0, fmt.Errorf("ожидалось число больше 0. Имеется %s: %w", songIdStr, err)
	}
	return idSong, nil
}

func validateBodyPatchSongRequest(req *EnrichedSong) error {
	const (
		layout = "02.01.2006"
	)

	var bufErrors bytes.Buffer

	parsedDate, err := time.Parse(layout, req.ReleaseDate)
	if err != nil && req.ReleaseDate != "" {
		bufErrors.WriteString(fmt.Sprintf("|дата не соответствует шаблону. Имеется: `%s`, должно быть: `%s`|",
			req.ReleaseDate, layout))
	}
	now := time.Now()
	if err == nil && parsedDate.After(now) && req.ReleaseDate != "" {
		bufErrors.WriteString(fmt.Sprintf("|дата песни не должна быть в будущем. Имеется: `%s`, сейчас: `%s`|",
			req.ReleaseDate, now.Format(layout)))
	}

	if _, err = url.Parse(req.Link); err != nil && req.Link != "" {
		bufErrors.WriteString(fmt.Sprintf("|нет ссылки на песню. Имеется: `%s`|",
			req.Link))
	}

	const (
		song         = "song"
		group        = "group"
		maxSymbSong  = 200
		maxSymbGroup = 100
	)

	symbGroup := []byte(req.Group)
	if req.Group != "" && len(symbGroup) > maxSymbGroup {
		bufErrors.WriteString(
			fmt.Sprintf("|для поля `%s` много символов: %d|", group, symbGroup))
	}
	symbSong := []byte(req.Song)
	if req.Song != "" && len(symbSong) > maxSymbSong {
		bufErrors.WriteString(
			fmt.Sprintf("|для поля `%s` много символов: %d|", song, symbSong))
	}

	if bufErrors.Len() > 0 {
		return errors.New(bufErrors.String())
	}

	return nil
}

// func (db *DataBase) Update(song *service.EnrichedSong, songID int, requestID string) (statusCode int, err error) {
func (s *Service) updSongInStorage(req *EnrichedSong, songID int, requestID string) ([]byte, int) {
	statusCode, err := s.delUpdater.Update(req, songID, requestID)
	if err != nil {
		dataJson, statusCode := createPatchDelSongResponse(
			false, statusCode, fmt.Sprint(err), requestID)
		return dataJson, statusCode
	}
	dataJson, statusCode := createPatchDelSongResponse(
		true, statusCode, msg200Upd, requestID)
	return dataJson, statusCode
}
