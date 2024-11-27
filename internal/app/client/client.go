package client

import (
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/service"
	"jedyEvgeny/online-music-library/pkg/logger"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	host           string
	path           string
	client         http.Client
	scheme         string
	maxRetriesConn int
	delayConn      int
	log            *logger.Logger
}

const (
	logStart           = "[%s] начинаем связь со сторонним сервером"
	logCreateRequest   = "[%s] сформирован запрос: %s"
	logErrRequest      = "[%s]не смогли создать запрос с URL: %s: %v"
	logErrExectRequest = "[%s] ошибка выполнения запроса №%d: %v"
	logRespStatus      = "[%s] получен код %d от стороннего сервера"
	logErrConnect      = "[%s] не удалось установить связь со сторонним сервером"
)

func New(host, port string, logger *logger.Logger, path string) *Client {
	return &Client{
		log:            logger,
		host:           createHost(host, port),
		path:           path,
		client:         http.Client{},
		scheme:         "http",
		maxRetriesConn: 3,
		delayConn:      1,
	}
}

func createHost(h, p string) string {
	return h + ":" + p
}

func (c *Client) Update(s *service.Song, requestID string) (*http.Response, error) {
	c.log.Debug(fmt.Sprintf(logStart, requestID))
	u := url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.path,
	}

	query := url.Values{}
	query.Add("group", s.Group)
	query.Add("song", s.Song)

	u.RawQuery = query.Encode()
	c.log.Debug(fmt.Sprintf(logCreateRequest, requestID, u.String()))

	var err error
	for attemptConn := 1; attemptConn <= c.maxRetriesConn; attemptConn++ {
		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			c.log.Debug(fmt.Sprintf(logErrRequest, requestID, u.String(), err))
			return nil, fmt.Errorf("не смогли создать запрос с URL: %s: %w",
				u.String(), err)
		}

		resp, err := c.client.Do(req)
		if err != nil {
			c.log.Debug(fmt.Sprintf(logErrExectRequest, requestID, attemptConn, err))

			delay := attemptConn * c.delayConn
			time.Sleep(time.Duration(delay) * time.Second)
			continue
		}

		c.log.Debug(fmt.Sprintf(logRespStatus, requestID, resp.StatusCode))

		switch resp.StatusCode {
		case http.StatusOK:
			return resp, nil
		case http.StatusNotFound:
			return nil, fmt.Errorf("ресурс в стороннем хранилище не найден, код ответа: %d", resp.StatusCode)
		case http.StatusInternalServerError:
			return nil, fmt.Errorf("внутренняя ошибка сервера, код ответа: %d", resp.StatusCode)
		default:
			return nil, fmt.Errorf("неизвестный ответ сервера: %s", resp.Status)
		}
	}

	c.log.Debug(fmt.Sprintf(logErrConnect, requestID))
	return nil, fmt.Errorf("не удалось связаться с сервером: %v", err)
}
