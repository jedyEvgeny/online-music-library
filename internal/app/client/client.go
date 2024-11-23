package client

import (
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/service"
	"log"
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
}

func New(host, port, path string) *Client {
	return &Client{
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

func (c *Client) Update(s *service.Song) (*http.Response, error) {
	u := url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   c.path,
	}

	query := url.Values{}
	query.Add("group", s.Group)
	query.Add("song", s.Song)

	u.RawQuery = query.Encode()

	var err error
	for attemptConn := 1; attemptConn <= c.maxRetriesConn; attemptConn++ {
		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("не смогли создать запрос с URL: %s: %w",
				u.String(), err)
		}
		log.Printf("сформирован запрос: %s\n", u.String())

		resp, err := c.client.Do(req)
		if err != nil {
			log.Printf("Ошибка выполнения запроса: %v", err)

			delay := attemptConn * c.delayConn
			time.Sleep(time.Duration(delay) * time.Second)
			continue
		}

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
	return nil, fmt.Errorf("не удалось связаться с сервером: %v", err)
}
