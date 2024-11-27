package service

import (
	"jedyEvgeny/online-music-library/pkg/logger"
	"net/http"
)

type WriteReader interface {
	Write(*EnrichedSong, string) (int, error)
	ReadLirycs(int, string) (string, int, error)
	ReadLibrary(*FilterAndPaggination, string) (*[]EnrichedSong, int, error)
}

type DelUpdater interface {
	Delete(int, string) error
	Update(*EnrichedSong, int, string) (int, error)
}

type Enricher interface {
	Update(*Song, string) (*http.Response, error)
}

type Service struct {
	log        *logger.Logger
	repository WriteReader
	enricher   Enricher
	delUpdater DelUpdater
}

func New(l *logger.Logger, w WriteReader, e Enricher, d DelUpdater) *Service {
	return &Service{
		log:        l,
		repository: w,
		enricher:   e,
		delUpdater: d,
	}
}
