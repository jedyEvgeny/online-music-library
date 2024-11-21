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

const (
	msg200 = "Ресурс существует"
	msg201 = "Ресурс создан"
)

const (
	errMarshalJson    = "ошибка создания json-объекта"
	errDecodeJson     = "ошибка декодирования json-объекта"
	errMethod         = "ошибка метода. Ожидался: %s, имеется: %s"
	errIsNotUUID      = "поле json valletID ожидалось с уникальным UUID. Имеется: %v"
	errOperation      = "поле json operationType: %s. Ожидалось 'DEPOSIT' или 'WITHDRAW'"
	errAmount         = "поле json amount должно быть больше нуля. Имеется: %d"
	errIsNotUUIDInURL = "в URL ожидалось UUID. Имеется: %s"
)
