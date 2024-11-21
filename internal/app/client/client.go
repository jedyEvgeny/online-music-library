package client

import (
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/service"
	"net/http"
	"net/url"
)

type Client struct {
	host   string
	path   string
	client http.Client
	scheme string
}

func New(host, path string) *Client {
	return &Client{
		host:   host,
		path:   path,
		client: http.Client{},
		scheme: "http",
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

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("не смогли создать запрос с URL: %s: %w",
			u.String(), err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("не удалось связаться с сервером: %w", err)
	}

	return resp, nil
}
