package storage

import (
	"database/sql"
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/service"
	"log"
	"time"
)

func (db *DataBase) Write(song *service.EnrichedSong, requestID string) (songID int, err error) {
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

	groupID, err := db.getOrCreateGroup(tx, song.Group)
	if err != nil {
		return 0, err
	}

	songID, err = db.findOrInsertSongID(tx, groupID, song, requestID)
	if err != nil {
		return 0, err
	}

	return songID, nil
}

func (db *DataBase) getOrCreateGroup(tx *sql.Tx, group string) (int, error) {
	var groupID int
	sqlStmt, err := tx.Prepare(requestSelectGroupID())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(group).Scan(&groupID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf(errExec, err)
	}
	if err == sql.ErrNoRows {
		return db.createGroup(tx, group)
	}
	return groupID, nil
}

func (db *DataBase) findOrInsertSongID(tx *sql.Tx, groupID int, song *service.EnrichedSong, requestID string) (int, error) {
	var songID int
	existingSongID, err := db.findSongIdByGroupID(tx, groupID, song.Song)
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
		err = db.updateSongToInitState(tx, existingSongID, song)
		if err != nil {
			return 0, err
		}
		songID = existingSongID
	}
	return songID, nil
}

func (db *DataBase) createGroup(tx *sql.Tx, group string) (int, error) {
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
	return groupID, nil
}

func (db *DataBase) findSongIdByGroupID(tx *sql.Tx, groupID int, songName string) (int, error) {
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
	return songID, nil
}

func (db *DataBase) updateSongToInitState(tx *sql.Tx, id int, song *service.EnrichedSong) error {
	sqlStmt, err := tx.Prepare(requestUpdateSongToinitState())
	if err != nil {
		return fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	_, err = sqlStmt.Exec(song.ReleaseDate, song.Lyrics, song.Link, id)
	if err != nil {
		return fmt.Errorf(errExec, err)
	}

	return nil
}

func (db *DataBase) addSong(tx *sql.Tx, groupID int, song *service.EnrichedSong, requestID string) (int, error) {
	var songID int

	sqlStmt, err := tx.Prepare(requestInsertSong())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	startInsert := time.Now()
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
	log.Printf(msgTimeInsert, requestID, time.Since(startInsert))

	return songID, nil
}
