package storage

import (
	"bytes"
	"database/sql"
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/service"
	"net/http"
	"time"
)

func (db *DataBase) Update(song *service.EnrichedSong, songID int, requestID string) (statusCode int, err error) {
	const op = "Update"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))
	startTx := time.Now()

	tx, err := db.db.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf(errCreateTx, err)
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

	err = db.checkSong(tx, song, songID, requestID)
	if err == sql.ErrNoRows {
		return http.StatusNotFound, err
	}
	if err != nil {
		return http.StatusInternalServerError, err
	}
	endTx := time.Now()

	db.log.Debug(fmt.Sprintf(logTxTime, requestID, endTx.Sub(startTx)))
	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return http.StatusOK, nil
}

func (db *DataBase) checkSong(tx *sql.Tx, e *service.EnrichedSong, songID int, requestID string) error {
	const op = "checkSong"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	currGroupID, err := db.findGroupID(tx, songID, requestID)
	if err != nil {
		return err
	}

	groupIdLikeGroupInNewSong, err := db.findExistGroupIdLikeNewGroup(tx, e.Group, requestID)
	if err != nil {
		return fmt.Errorf(errNoContent, err)
	}

	foundMoreThanOneSongForCurrGroup, err := db.hasMoreThanOneSongForCurrGroup(tx, currGroupID, requestID)
	if err != nil {
		return fmt.Errorf(errNoContent, err)
	}

	if e.Group != "" {
		err = db.updGroup(tx, e, foundMoreThanOneSongForCurrGroup, songID, groupIdLikeGroupInNewSong, requestID)
		if err != nil {
			return err
		}
	}

	needToChange, err := db.updSong(tx, e, songID, requestID)
	if !needToChange && err == nil { //если поменялось только название исполнителя
		return nil
	}
	if err != nil {
		return err
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return nil
}

// findGroupID возвращает id группы из таблицы songs
// и косвенно проверяет, существует ли переданный songID
func (db *DataBase) findGroupID(tx *sql.Tx, songID int, requestID string) (int, error) {
	const op = "findGroupID"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	var currGroupID int

	sqlStmt, err := tx.Prepare(requestSelectGroupIdBySongID())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(songID).Scan(&currGroupID)
	if err != nil {
		return 0, fmt.Errorf(errExec, err)
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return currGroupID, nil
}

func (db *DataBase) findExistGroupIdLikeNewGroup(tx *sql.Tx, group, requestID string) (int, error) {
	const op = "findExistGroupIdLikeNewGroup"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	var groupIdLikeGroupInNewSong int

	sqlStmt, err := tx.Prepare(requestSelectGroupIdByGroupName())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(group).Scan(&groupIdLikeGroupInNewSong)
	if err == sql.ErrNoRows {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf(errExec, err)
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return groupIdLikeGroupInNewSong, nil
}

func (db *DataBase) hasMoreThanOneSongForCurrGroup(tx *sql.Tx, groupID int, requestID string) (bool, error) {
	const op = "hasMoreThanOneSongForCurrGroup"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	var songCount int

	sqlStmt, err := tx.Prepare(requestSelectCountSongsByGroup())
	if err != nil {
		return false, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(groupID).Scan(&songCount)
	if err != nil {
		return false, fmt.Errorf(errExec, err)
	}

	if songCount == 1 {
		return false, nil
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return true, nil
}

func (db *DataBase) updGroup(tx *sql.Tx, e *service.EnrichedSong, foundMoreTwoSongsByCurrGroup bool, songID, groupIdLikeGroupInNewSong int, requestID string) error {
	const op = "updGroup"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	if (!foundMoreTwoSongsByCurrGroup && groupIdLikeGroupInNewSong == 0) ||
		(foundMoreTwoSongsByCurrGroup && groupIdLikeGroupInNewSong == 0) {
		newGroupID, err := createNewGroup(tx, e)
		if err != nil {
			return err
		}

		err = updGroupIdInSongsBySongID(tx, newGroupID, songID)
		if err != nil {
			return err
		}
	}

	if (!foundMoreTwoSongsByCurrGroup && groupIdLikeGroupInNewSong != 0) ||
		(foundMoreTwoSongsByCurrGroup && groupIdLikeGroupInNewSong != 0) {
		err := updGroupIdInSongsBySongID(tx, groupIdLikeGroupInNewSong, songID)
		if err != nil {
			return err
		}
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return nil
}

func createNewGroup(tx *sql.Tx, e *service.EnrichedSong) (int, error) {
	sqlStmt, err := tx.Prepare(requestInsertGroup())
	if err != nil {
		return 0, fmt.Errorf(errExec, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	var newGroupID int

	err = sqlStmt.QueryRow(e.Group).Scan(&newGroupID)
	if err != nil {
		return 0, fmt.Errorf(errExec, err)
	}

	return newGroupID, nil
}

func updGroupIdInSongsBySongID(tx *sql.Tx, newGroupID, songID int) error {
	sqlStmt, err := tx.Prepare(requestUpdateGroupIdInSongsBySongID())
	if err != nil {
		return fmt.Errorf(errExec, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	_, err = sqlStmt.Exec(
		newGroupID,
		songID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBase) updSong(tx *sql.Tx, e *service.EnrichedSong, songID int, requestID string) (bool, error) {
	const op = "updSong"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	query, args := db.createQueryUpdSong(e, songID, requestID)

	if len(args) == 0 {
		return false, nil
	}

	sqlStmt, err := tx.Prepare(query)
	if err != nil {
		return false, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	_, err = sqlStmt.Exec(args...)
	if err != nil {
		return false, fmt.Errorf(errExec, err)
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return true, nil
}

func (db *DataBase) createQueryUpdSong(e *service.EnrichedSong, songID int, requestID string) (string, []interface{}) {
	const op = "createQueryUpdSong"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	var buf bytes.Buffer
	var args []interface{}

	buf.WriteString("UPDATE songs SET ")

	if e.Song != "" {
		args = append(args, e.Song)
		buf.WriteString(fmt.Sprintf("song = $%d, ", len(args)))
	}
	if e.ReleaseDate != "" {
		args = append(args, e.ReleaseDate)
		buf.WriteString(fmt.Sprintf("release_date = $%d, ", len(args)))
	}
	if e.Lyrics != "" {
		args = append(args, e.Lyrics)
		buf.WriteString(fmt.Sprintf("lyrics = $%d, ", len(args)))
	}
	if e.Link != "" {
		args = append(args, e.Link)
		buf.WriteString(fmt.Sprintf("link = $%d, ", len(args)))
	}

	if len(args) > 0 {
		args = append(args, songID)
		buf.Truncate(buf.Len() - 2) //убираем последнюю зпт из запроса
		buf.WriteString(fmt.Sprintf(" WHERE s_id = $%d", len(args)))
	}
	query := buf.String()

	db.log.Debug(fmt.Sprintf("sql-запрос: %s", query))
	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return query, args
}
