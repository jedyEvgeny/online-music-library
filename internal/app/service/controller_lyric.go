package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func (s *Service) ProcessReadLirycsSongRequest(r *http.Request, requestID string) ([]byte, int) {
	defer closeRequestBody(r.Body)

	reqParam, statusCode, errResponse := validateReadLirycsRequest(r, requestID)
	if statusCode != http.StatusOK {
		return errResponse, statusCode
	}

	liryc, statusCode, err := s.repository.ReadLirycs(reqParam.idSong, requestID)
	if err != nil {
		dataJson, statusCode := createLirycsResponse(
			false, statusCode, fmt.Sprint(err), requestID, nil)
		return dataJson, statusCode
	}

	lirycs := createSliceLirycs(liryc, reqParam)

	dataJson, statusCode := createLirycsResponse(
		true, http.StatusOK, msg200, requestID, &lirycs)
	return dataJson, statusCode
}

func validateReadLirycsRequest(r *http.Request, requestID string) (*paggination, int, []byte) {
	if r.Method != http.MethodGet {
		msg := fmt.Sprintf(errMethod, http.MethodGet, r.Method)
		dataJson, statusCode := createAddSongResponse(
			false, http.StatusMethodNotAllowed, msg, requestID, nil)
		return nil, statusCode, dataJson
	}

	param := decodeLirycsRequest(r.URL)

	err := param.validateLyricsParamsRequest()
	if err != nil {
		dataJson, statusCode := createLirycsResponse(
			false, http.StatusBadRequest, fmt.Sprint(err), requestID, nil)
		return nil, statusCode, dataJson
	}

	param.fillFields()

	return param, http.StatusOK, nil
}

func decodeLirycsRequest(u *url.URL) *paggination {
	path := u.Path
	log.Println("Путь:", path)

	q := u.Query()

	offset := q.Get("offset")
	if offset == "" {
		offset = "0"
	}

	limit := q.Get("limit")
	if limit == "" {
		limit = "0"
	}

	p := &paggination{
		offsetStr: offset,
		limitStr:  limit,
		idSongStr: strings.TrimPrefix(path, "/lyrics/"),
	}

	log.Println(*p)

	return p
}

func (p *paggination) validateLyricsParamsRequest() error {
	var buf bytes.Buffer

	idSong, err := strconv.Atoi(p.idSongStr)
	if err != nil {
		buf.WriteString(fmt.Sprintf("ошибка: %v. `id_song` не число. Имеется: %s", err, p.idSongStr))
	}
	if err == nil && idSong <= 0 {
		buf.WriteString(fmt.Sprintf("ошибка: %v. `id_song` не может быть меньше 1. Имеется: %s", err, p.idSongStr))
	}

	offset, err := strconv.Atoi(p.offsetStr)
	if err != nil {
		buf.WriteString(fmt.Sprintf("ошибка: %v. `offset` не число. Имеется: %s", err, p.offsetStr))
	}
	if err == nil && offset < 0 {
		buf.WriteString(fmt.Sprintf("ошибка: %v. `offset` не может быть меньше 1. Имеется: %s", err, p.offsetStr))
	}

	limit, err := strconv.Atoi(p.limitStr)
	if err != nil {
		buf.WriteString(fmt.Sprintf("ошибка: %v. `limit` не число. Имеется: %s", err, p.limitStr))
	}
	if err == nil && limit < 0 {
		buf.WriteString(fmt.Sprintf("ошибка: %v. `limit` не может быть меньше 1. Имеется: %s", err, p.limitStr))
	}

	if buf.Len() > 0 {
		return errors.New(buf.String())
	}

	return nil
}

func createLirycsResponse(ok bool, statusCode int, msg, requestID string, lirycs *[]string) ([]byte, int) {
	log.Printf("[%s]  %s\n", requestID, msg)
	resp := ResponseLirycs{
		Sucsess:    ok,
		Message:    msg,
		StatusCode: statusCode,
		Lirycs:     *lirycs,
	}
	dataJson, err := json.Marshal(resp)
	if err != nil {
		log.Printf(errMarshalJson, err)
		return nil, http.StatusInternalServerError
	}
	return dataJson, statusCode
}

func (p *paggination) fillFields() {
	p.idSong, _ = strconv.Atoi(p.idSongStr)
	p.offset, _ = strconv.Atoi(p.offsetStr)
	p.limit, _ = strconv.Atoi(p.limitStr)
}

func createSliceLirycs(s string, p *paggination) []string {
	lirycs := strings.Split(s, "\\n\\n")

	if p.offset >= len(lirycs) {
		return nil
	}

	endIdx := p.offset + p.limit
	if endIdx > len(lirycs) {
		endIdx = len(lirycs)
	}

	return lirycs[p.offset:endIdx]
}
