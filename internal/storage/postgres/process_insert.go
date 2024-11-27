package storage

import (
	"database/sql"
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/service"
	"time"
)

func (db *DataBase) Write(song *service.EnrichedSong, requestID string) (songID int, err error) {
	const op = "Write"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	startTx := time.Now()

	tx, err := db.db.Begin()
	if err != nil {
		return 0, fmt.Errorf(errCreateTx, err)
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

	groupID, err := db.getOrCreateGroup(tx, song.Group, requestID)
	if err != nil {
		return 0, err
	}

	songID, err = db.findOrInsertSongID(tx, groupID, song, requestID)
	if err != nil {
		return 0, err
	}
	endTx := time.Now()

	db.log.Debug(fmt.Sprintf(logTxTime, requestID, endTx.Sub(startTx)))
	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return songID, nil
}

func (db *DataBase) getOrCreateGroup(tx *sql.Tx, group, requestID string) (int, error) {
	const op = "getOrCreateGroup"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	var groupID int
	sqlStmt, err := tx.Prepare(requestSelectGroupIdByGroupName())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(group).Scan(&groupID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf(errExec, err)
	}
	if err == sql.ErrNoRows {
		return db.createGroup(tx, group, requestID)
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return groupID, nil
}

func (db *DataBase) findOrInsertSongID(tx *sql.Tx, groupID int, song *service.EnrichedSong, requestID string) (int, error) {
	const op = "findOrInsertSongID"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	var songID int
	existingSongID, err := db.findSongIdByGroupID(tx, groupID, song.Song, requestID)
	if err != nil {
		return 0, err
	}
	if existingSongID == 0 {
		songID, err = db.addSong(tx, groupID, song, requestID)
		if err != nil {
			return 0, err
		}
	}
	if existingSongID != 0 {
		err = db.updateSongToInitState(tx, existingSongID, song, requestID)
		if err != nil {
			return 0, err
		}
		songID = existingSongID
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return songID, nil
}

func (db *DataBase) createGroup(tx *sql.Tx, group, requestID string) (int, error) {
	const op = "createGroup"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	sqlStmt, err := tx.Prepare(requestInsertGroup())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	var groupID int
	err = sqlStmt.QueryRow(group).Scan(&groupID)
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return groupID, nil
}

func (db *DataBase) findSongIdByGroupID(tx *sql.Tx, groupID int, songName, requestID string) (int, error) {
	const op = "findSongIdByGroupID"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	var songID int
	sqlStmt, err := tx.Prepare(requestSelectSongByGroup())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(
		groupID,
		songName,
	).Scan(&songID)

	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf(errExec, err)
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return songID, nil
}

func (db *DataBase) updateSongToInitState(tx *sql.Tx, id int, song *service.EnrichedSong, requestID string) error {
	const op = "updateSongToInitState"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	sqlStmt, err := tx.Prepare(requestUpdateSongToinitState())
	if err != nil {
		return fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	_, err = sqlStmt.Exec(song.ReleaseDate, song.Lyrics, song.Link, id)
	if err != nil {
		return fmt.Errorf(errExec, err)
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return nil
}

func (db *DataBase) addSong(tx *sql.Tx, groupID int, song *service.EnrichedSong, requestID string) (int, error) {
	const op = "addSong"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	var songID int

	sqlStmt, err := tx.Prepare(requestInsertSong())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(
		groupID,
		song.Song,
		song.ReleaseDate,
		song.Lyrics,
		song.Link,
	).Scan(&songID)
	if err != nil {
		return 0, fmt.Errorf(errExec, err)
	}

	db.log.Debug(fmt.Sprintf(logStart, requestID, op))
	return songID, nil
}
