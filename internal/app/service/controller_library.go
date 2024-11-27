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

const (
	Group       = "group"
	Composition = "song"
	ReleaseDate = "releaseDate"
)

func (s *Service) ProcessLibraryRequest(r *http.Request, requestID string) ([]byte, int) {
	defer closeRequestBody(r.Body)

	reqParam, statusCode, errResponse := s.validateLibraryRequest(r, requestID)
	if statusCode != http.StatusOK {
		s.log.Debug(fmt.Sprintf(logErrValidate, requestID, string(errResponse)))
		return errResponse, statusCode
	}

	songs, statusCode, err := s.repository.ReadLibrary(reqParam, requestID)
	s.log.Debug(fmt.Sprintf(logAnswDB, statusCode, err))
	if err != nil {
		dataJson, statusCode := s.createLibraryResponse(
			false, statusCode, fmt.Sprint(err), requestID, nil)
		return dataJson, statusCode
	}

	dataJson, statusCode := s.createLibraryResponse(
		true, http.StatusOK, msg200, requestID, songs)

	s.log.Debug(fmt.Sprintf(logToEndpoin, requestID, statusCode))
	return dataJson, statusCode
}

func (s *Service) validateLibraryRequest(r *http.Request, requestID string) (*FilterAndPaggination, int, []byte) {
	if r.Method != http.MethodGet {
		msg := fmt.Sprintf(errMethod, http.MethodGet, r.Method)
		dataJson, statusCode := s.createLibraryResponse(
			false, http.StatusMethodNotAllowed, msg, requestID, nil)
		return nil, statusCode, dataJson
	}

	param := decodeLibraryRequest(r.URL)

	err := param.validateAndFillLibraryParamsRequest()
	if err != nil {
		dataJson, statusCode := s.createLibraryResponse(
			false, http.StatusBadRequest, fmt.Sprint(err), requestID, nil)
		return nil, statusCode, dataJson
	}

	return param, http.StatusOK, nil
}

func decodeLibraryRequest(u *url.URL) *FilterAndPaggination {
	path := u.Path
	log.Println("Путь:", path)

	q := u.Query()

	sortBy := q.Get("sortBy")
	if sortBy == "" {
		sortBy = "asc"
	}

	offset := q.Get("offset")
	if offset == "" {
		offset = "0"
	}

	limit := q.Get("limit")
	if limit == "" {
		limit = "0"
	}

	p := &FilterAndPaggination{
		SortBy:    sortBy,
		filter:    q.Get("filter"),
		offsetStr: offset,
		limitStr:  limit,
	}

	log.Println(*p)
	return p
}

func (f *FilterAndPaggination) validateAndFillLibraryParamsRequest() error {
	var buf bytes.Buffer

	if f.SortBy != "asc" && f.SortBy != "desc" {
		buf.WriteString(fmt.Sprintf("параметр `sortBy` ожидалось `asc` или `desc`. Имеется: %s", f.SortBy))
	}

	var err error
	filterMap := make(map[string]interface{})
	if f.filter != "" {
		filterMap, err = validateAndReturnFilterMap(f.filter)
		if err != nil {
			buf.WriteString(fmt.Sprintf("ошибка: %vне валидный фильтр. Имеется: %s", err, f.filter))
		}
	}

	offset, err := strconv.Atoi(f.offsetStr)
	if err != nil {
		buf.WriteString(fmt.Sprintf("ошибка: %v. `offset` не число. Имеется: %s", err, f.offsetStr))
	}
	if err == nil && offset < 0 {
		buf.WriteString(fmt.Sprintf("ошибка: %v. `offset` не может быть меньше 1. Имеется: %s", err, f.offsetStr))
	}

	limit, err := strconv.Atoi(f.limitStr)
	if err != nil {
		buf.WriteString(fmt.Sprintf("ошибка: %v. `limit` не число. Имеется: %s", err, f.limitStr))
	}
	if err == nil && limit < 0 {
		buf.WriteString(fmt.Sprintf("ошибка: %v. `limit` не может быть меньше 1. Имеется: %s", err, f.limitStr))
	}

	if buf.Len() > 0 {
		return errors.New(buf.String())
	}

	f.Limit = limit
	f.Offset = offset
	f.Filter = filterMap

	log.Printf("Тип сортировки: %s | offset=%d | limit=%d | фильтр=%s",
		f.SortBy, f.Offset, f.Limit, f.Filter)
	return nil
}

func validateAndReturnFilterMap(filter string) (map[string]interface{}, error) {
	splits := strings.SplitAfterN(filter, ".", 2)
	splits[0] = strings.TrimSuffix(splits[0], ".")
	if len(splits) != 2 {
		return nil, fmt.Errorf("невалидный разделитель фильтра. Ожидался формат`поле.фильтр` с двумя элементами. Имеется количество элементов: %d", len(splits))
	}
	field, value := splits[0], splits[1]
	if field != Group && field != Composition && field != ReleaseDate {
		return nil, fmt.Errorf("невалидное поле фильтра. Ожидалось `%s/%s/%s`. Имеется: %s",
			Group, Composition, ReleaseDate, field)
	}

	return map[string]interface{}{field: value}, nil
}

func (s *Service) createLibraryResponse(ok bool, statusCode int, msg, requestID string, library *[]EnrichedSong) ([]byte, int) {
	s.log.Debug(fmt.Sprintf("[%s]  %s\n", requestID, msg))
	resp := ResponseLibrary{
		Sucsess:    ok,
		Message:    msg,
		StatusCode: statusCode,
	}
	if library != nil {
		resp.Library = *library
	}
	dataJson, err := json.Marshal(resp)
	if err != nil {
		log.Printf(errMarshalJson, err)
		return nil, http.StatusInternalServerError
	}
	return dataJson, statusCode
}
