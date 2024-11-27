package service

type Song struct {
	Group string `json:"group"`
	Song  string `json:"song"`
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

type ResponsePost struct {
	Sucsess    bool   `json:"sucsess"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	ID         *int   `json:"resourceID,omitempty"`
}

type EnrichedSong struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Lyrics      string `json:"text"`
	Link        string `json:"link"`
}

type ResponsePatchDelete struct {
	Sucsess    bool   `json:"sucsess"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

type ResponseLirycs struct {
	Sucsess    bool     `json:"sucsess"`
	Message    string   `json:"message"`
	StatusCode int      `json:"statusCode"`
	Lirycs     []string `json:"lyrics,omitempty"`
}

type ResponseLibrary struct {
	Sucsess    bool           `json:"sucsess"`
	Message    string         `json:"message"`
	StatusCode int            `json:"statusCode"`
	Library    []EnrichedSong `json:"library,omitempty"`
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
	msg200    = "Ресурс существует"
	msg201    = "Ресурс создан"
	msg204    = "Ресурс в хранилище отсутвует"
	msg200Upd = "Ресурс обновлён"
)

const (
	logErrValidate = "[%s] запрос не прошёл валидацию: %s"
	logToEndpoin   = "[%s] подготовлен ответ со статусом: %d"
	logAnswDB      = "Ответ БД. Статус: %d, ошибка: %v\n"
)

const (
	errMarshalJson = "ошибка создания json-объекта: %v"
	errDecodeJson  = "ошибка декодирования json-объекта"
	errMethod      = "ошибка метода. Ожидался: %s, имеется: %s"
	errIDDel       = "не смогли прочитать параметр `s_id` в строке запроса: %w"
	errIDValueDel  = "`s_id` не должно быть меньше 1. Имеется: %d"
)
