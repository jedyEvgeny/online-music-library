package service

import "time"

type Song struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type ResponsePost struct {
	Sucsess    bool   `json:"sucsess"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	ID         *int   `json:"resourceID,omitempty"`
}

type EnrichedSong struct {
	Group           string `json:"group"`
	Song            string `json:"song"`
	ReleaseDate     string `json:"releaseDate"`
	ReleaseDateTime time.Time
	Lyrics          string `json:"text"`
	Link            string `json:"link"`
}

type ResponseDelete struct {
	Sucsess    bool   `json:"sucsess"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type ResponseLirycs struct {
	Sucsess    bool     `json:"sucsess"`
	Message    string   `json:"message"`
	StatusCode int      `json:"statusCode"`
	Lirycs     []string `json:"lirycs,omitempty"`
}

type ResponseLibrary struct {
	Sucsess    bool           `json:"sucsess"`
	Message    string         `json:"message"`
	StatusCode int            `json:"statusCode"`
	Library    []EnrichedSong `json:"library,omitempty"`
}

type FilterAndPaggination struct {
	SortBy    string
	Filter    map[string]interface{}
	Offset    int
	Limit     int
	filter    string
	offsetStr string
	limitStr  string
}

type paggination struct {
	offsetStr string
	limitStr  string
	idSongStr string
	offset    int
	limit     int
	idSong    int
}

const (
	msg200 = "Ресурс существует"
	msg201 = "Ресурс создан"
	msg204 = "Ресурс в хранилище отсутвует"
)

const (
	errMarshalJson = "ошибка создания json-объекта: %v"
	errDecodeJson  = "ошибка декодирования json-объекта"
	errMethod      = "ошибка метода. Ожидался: %s, имеется: %s"
	errIsNotUUID   = "поле json valletID ожидалось с уникальным UUID. Имеется: %v"
	errOperation   = "поле json operationType: %s. Ожидалось 'DEPOSIT' или 'WITHDRAW'"
	errAmount      = "поле json amount должно быть больше нуля. Имеется: %d"
	errIDDel       = "не смогли прочитать параметр `s_id` в строке запроса: %w"
	errIDValueDel  = "`s_id` не должно быть меньше 1. Имеется: %d"
)
