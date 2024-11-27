package endpoint

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

type delAdder interface {
	ProseccAddSongRequest(*http.Request, string) ([]byte, int)
	ProseccDelSongRequest(*http.Request, string) ([]byte, int)
}

type UpdateReader interface {
	ProcessUpdateSongRequest(*http.Request, string) ([]byte, int)
	ProcessReadLirycsSongRequest(*http.Request, string) ([]byte, int)
	ProcessLibraryRequest(*http.Request, string) ([]byte, int)
}

type Endpoint struct {
	process delAdder
	update  UpdateReader
}

func New(c delAdder, u UpdateReader) *Endpoint {
	return &Endpoint{
		process: c,
		update:  u,
	}
}

const (
	msgRequest = "[%s] Получен запрос с методом: %s от URL: %s. Обработчик: %s\n"
)

// @Summary		Добавить песню
// @Description	Добавляет новую песню.
// @Description	Наименование песни и группа передаются в теле запроса в json-объекте.
// @Description	При создании песни происходит обращение к удалённому серверу для обогащения информации.
// @Tags			songs
// @Accept			json
// @Produce		json
// @Param			song	body		Song			true	"Добавляем песню"
// @Success		201		{object}	ResponsePost201	"Запись успешно создана"
// @Failure		400		{object}	ResponsePost400	"Ошибка валидации данных"
// @Failure		405		{object}	ResponsePost405	"Метод не разрешен"
// @Failure		500		{object}	ResponsePost500	"Ошибка сервера"
// @Router			/song-add [post]
func (e *Endpoint) HandlerAddSong(w http.ResponseWriter, r *http.Request) {
	h := "HandlerAddSong"
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL, h)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resp, status := e.process.ProseccAddSongRequest(r, reqID)
	w.WriteHeader(status)
	w.Write(resp)
}

// @Summary		Удалить песню
// @Description	Удаляет песню.
// @Description	ID песни передаётся в URL.
// @Description	При отсутствии песни возвращается статус 204, как если бы песня была и успешно удалена.
// @Tags			songs
// @Produce		json
// @Param			songID	path	int	true	"id существующей песни"
// @Success		204		"Запись успешно создана"
// @Failure		400		{object}	ResponseDel400	"Ошибка валидации данных"
// @Failure		405		{object}	ResponseDel405	"Метод не разрешен"
// @Failure		500		{object}	ResponsePost500	"Ошибка сервера"
// @Router			/song-del/{songID} [delete]
func (e *Endpoint) HandlerDeleteSong(w http.ResponseWriter, r *http.Request) {
	h := "HandlerDeleteSong"
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL, h)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, status := e.process.ProseccDelSongRequest(r, reqID)
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

// @Summary		Получить текст песни
// @Description	Получает текст песни с пагинацией по куплетам.
// @Description	ID песни передаётся в URL, параметры пагинации в параметрах запроса.
// @Tags			songs
// @Produce		json
// @Param			songID	path		int					true	"id существующей песни"
// @Param			offset	query		int					false	"Смещение для пагинации"
// @Param			limit	query		int					false	"Количество записей для пагинации"
// @Success		200		{object}	ResponseLirycs200	"Запись успешно создана"
// @Failure		400		{object}	ResponseLirycs400	"Ошибка валидации данных"
// @Failure		404		{object}	ResponseLirycs404	"Ресурс не найден"
// @Failure		405		{object}	ResponseLirycs405	"Метод не разрешен"
// @Failure		500		{object}	ResponseLirycs500	"Ошибка сервера"
// @Router			/lyrics/{songID} [get]
func (e *Endpoint) HandlerLiryc(w http.ResponseWriter, r *http.Request) {
	h := "HandlerLiryc"
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL, h)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, status := e.update.ProcessReadLirycsSongRequest(r, reqID)
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

// @Summary		Получить перечень песен
// @Description	Получает данные библиотеки с фильтрацией и пагинацией.
// @Description	Параметры фильтрации и пагинации передаются в параметрах запроса. ID песни передаётся как часть пути в URL.
// @Description	Фильтр даты передавать в формате
// @Tags			songs
// @Produce		json
// @Param			songID	path		int				true	"id существующей песни"
// @Param			filter	query		string				true	"Фильтр поиска. Возможные значения: releaseDate, group, song, например releaseDate.26.08.1968, group.Muse, song.Supermassive Black Hole"
// @Param			offset	query		int					false	"Смещение для пагинации"
// @Param			limit	query		int					false	"Количество записей для пагинации"
// @Success		200		{object}	ResponseLibrary200	"Запись успешно создана"
// @Failure		400		{object}	ResponseLibrary400	"Ошибка валидации данных"
// @Failure		404		{object}	ResponseLibrary404	"Ресурс не найден"
// @Failure		405		{object}	ResponseLibrary405	"Метод не разрешен"
// @Failure		500		{object}	ResponseLibrary500	"Ошибка сервера"
// @Router			/list/{songID} [get]
func (e *Endpoint) HandlerLibrary(w http.ResponseWriter, r *http.Request) {
	h := "HandlerLibrary"
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL, h)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, status := e.update.ProcessLibraryRequest(r, reqID)
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

// @Summary		Изменить параметры песни
// @Description	Изменяет один или несколько параметров песни.
// @Description	ID песни передаётся как часть пути, параметры песни - в теле запроса.
// @Tags			songs
// @Accept			json
// @Produce		json
// @Param			song	body		EnrichedSong			true	"Возможные поля для изменения"
// @Param			songID	path		int				true	"id существующей песни"
// @Success		200		{object}	ResponseUpdate200	"Запись успешно создана"
// @Failure		400		{object}	ResponseLibrary400	"Ошибка валидации данных"
// @Failure		404		{object}	ResponseLibrary404	"Ресурс не найден"
// @Failure		405		{object}	ResponseUpdate405	"Метод не разрешен"
// @Failure		500		{object}	ResponseLibrary500	"Ошибка сервера"
// @Router			/song-upd/{songID} [patch]
func (e *Endpoint) HandlerPatchSong(w http.ResponseWriter, r *http.Request) {
	h := "HandlerPatchSong"
	reqID := requestID()
	log.Printf(msgRequest, reqID, r.Method, r.URL, h)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp, status := e.update.ProcessUpdateSongRequest(r, reqID)
	if status == http.StatusNoContent {
		w.WriteHeader(status)
		return
	}
	w.WriteHeader(status)
	w.Write(resp)
}

func requestID() string {
	return uuid.New().String()
}
