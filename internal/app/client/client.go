package client

import (
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/service"
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

func New(host, path string) *Client {
	return &Client{
		host:           host,
		path:           path,
		client:         http.Client{},
		scheme:         "http",
		maxRetriesConn: 3,
		delayConn:      1,
	}
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

	var resp *http.Response
	var err error
	for attemptConn := 1; attemptConn <= c.maxRetriesConn; attemptConn++ {
		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("не смогли создать запрос с URL: %s: %w",
				u.String(), err)
		}

		resp, err = c.client.Do(req)
		if err == nil {
			return resp, nil
		}
		delay := attemptConn * c.delayConn
		time.Sleep(time.Duration(delay) * time.Second)
	}
	return nil, fmt.Errorf("не удалось связаться с сервером: %v", err)
}
