package storage

import (
	"database/sql"
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/service"
	"net/http"
)

func (db *DataBase) ReadLibrary(f *service.FilterAndPaggination, requestID string) (songs *[]service.EnrichedSong, statusCode int, err error) {
	tx, err := db.db.Begin()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(errCreateTx, err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		commitErr := tx.Commit()
		if commitErr != nil {
			err = fmt.Errorf(errCommitTx, commitErr)
		}
	}()

	songs, statusCode, err = db.findSongsByFilter(tx, f, requestID)
	if err != nil && err != sql.ErrNoRows {
		return nil, http.StatusInternalServerError, err
	}
	if err == sql.ErrNoRows {
		return nil, http.StatusNotFound, err
	}

	return songs, http.StatusOK, nil
}

func (db *DataBase) ReadLirycs(songID int, requestID string) (liryc string, statusCode int, err error) {
	tx, err := db.db.Begin()
	if err != nil {
		return "", http.StatusInternalServerError, fmt.Errorf(errCreateTx, err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
			return
		}
		commitErr := tx.Commit()
		if commitErr != nil {
			err = fmt.Errorf(errCommitTx, commitErr)
		}
	}()

	liryc, err = db.findLirycsBySongID(tx, songID, requestID)
	if err != nil && err != sql.ErrNoRows {
		return "", http.StatusInternalServerError, err
	}
	if err == sql.ErrNoRows {
		return "", http.StatusNotFound, err
	}

	return liryc, 0, nil
}

func (db *DataBase) findLirycsBySongID(tx *sql.Tx, songID int, requestID string) (string, error) {
	var liryc string

	sqlStmt, err := tx.Prepare(requestSelectLirycsBySongID())
	if err != nil {
		return "", fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(
		songID,
	).Scan(&liryc)

	if err != nil {
		return "", fmt.Errorf(errExec, err)
	}
	return liryc, nil
}

func (db *DataBase) findSongsByFilter(tx *sql.Tx, f *service.FilterAndPaggination, requestID string) (*[]service.EnrichedSong, int, error) {
	query, arg := createQuery(f)

	sqlStmt, err := tx.Prepare(query)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	rows, err := sqlStmt.Query(arg, f.Limit, f.Offset)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(errExec, err)
	}
	defer rows.Close()

	var songs []service.EnrichedSong
	for rows.Next() {
		var song service.EnrichedSong
		err = rows.Scan(&song.Song, &song.ReleaseDate, &song.Lyrics, &song.Link, &song.Group)
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf(errRead, err)
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf(errRead, err)
	}
	if len(songs) == 0 {
		return nil, http.StatusNotFound, fmt.Errorf(errNoContent, err)
	}

	return &songs, http.StatusOK, nil
}

func createQuery(f *service.FilterAndPaggination) (string, interface{}) {
	var args interface{}
	var condition string

	if len(f.Filter) == 1 {
		for key, value := range f.Filter {
			prefix := "s"
			switch key {
			case service.Group:
				key = "name"
				prefix = "mg"
			case service.ReleaseDate:
				key = "release_date"
			}

			condition = fmt.Sprintf("%s.%s = $1", prefix, key)

			args = value
		}
	}

	if condition != "" {
		condition = fmt.Sprintf(" WHERE %s", condition)
	}

	// Доработать сортировку по ключу
	// orderBy := fmt.Sprintf(" ORDER BY s.song %s", f.SortBy)
	// query := requestSelectSongsByFilter(condition) + orderBy + " LIMIT $2 OFFSET $3"

	query := requestSelectSongsByFilter(fmt.Sprintf("%s  LIMIT $2 OFFSET $3", condition))

	return query, args
}
