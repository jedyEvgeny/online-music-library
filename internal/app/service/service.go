package service

import (
	"net/http"
)

type WriteReader interface {
	Write(*EnrichedSong, string) (int, error)
	ReadLirycs(int, string) (string, int, error)
	ReadLibrary(*FilterAndPaggination, string) (*[]EnrichedSong, int, error)
}

// func (db *DataBase) Update(song *service.EnrichedSong, songID int, requestID string) (statusCode int, err error) {

type DelUpdater interface {
	Delete(int, string) error
	Update(*EnrichedSong, int, string) (int, error)
}

type Enricher interface {
	Update(*Song) (*http.Response, error)
}

type Service struct {
	repository WriteReader
	enricher   Enricher
	delUpdater DelUpdater
}

func New(w WriteReader, e Enricher, d DelUpdater) *Service {
	return &Service{
		repository: w,
		enricher:   e,
		delUpdater: d,
	}
}
