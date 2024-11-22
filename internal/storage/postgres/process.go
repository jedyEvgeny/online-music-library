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

	groupID, err := db.getOrInsertGroup(tx, song.Group)
	if err != nil {
		return 0, err
	}
	songID, err = db.insertSong(tx, groupID, song, requestID)
	if err != nil {
		return 0, err
	}
	return songID, nil
}

func (db *DataBase) getOrInsertGroup(tx *sql.Tx, group string) (int, error) {
	var groupID int
	sqlStmt, err := tx.Prepare(reqSelectGroupID())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer sqlStmt.Close()

	err = sqlStmt.QueryRow(group).Scan(&groupID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf(errExec, err)
	}
	if err == sql.ErrNoRows {
		return db.insertGroup(tx, group)
	}
	return groupID, nil
}

func (db *DataBase) insertGroup(tx *sql.Tx, group string) (int, error) {
	sqlStmt, err := tx.Prepare(requestInsertGroup())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer sqlStmt.Close()

	var groupID int
	err = sqlStmt.QueryRow(group).Scan(&groupID)
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	return groupID, nil
}

func (db *DataBase) insertSong(tx *sql.Tx, groupID int, song *service.EnrichedSong, requestID string) (int, error) {
	var songID int

	sqlStmt, err := tx.Prepare(requestInsertSong())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer sqlStmt.Close()

	startInsert := time.Now()
	err = sqlStmt.QueryRow(
		groupID,
		song.Song,
		song.ReleaseDateTime,
		song.Lyrics,
		song.Link,
	).Scan(&songID)
	if err != nil {
		return 0, fmt.Errorf(errExec, err)
	}
	log.Printf(msgTimeInsert, requestID, time.Since(startInsert))

	return songID, nil
}
