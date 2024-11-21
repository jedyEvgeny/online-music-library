package storage

import (
	"database/sql"
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/service"
	"log"
	"time"
)

func (db *DataBase) Write(song *service.EnrichedSong, requestID string) (int, error) {
	var groupID, songID int

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
			err = fmt.Errorf(errCommitTx, err)
		}
	}()

	sqlStmt, err := tx.Prepare(reqSelectGroupID())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(song.Group).Scan(&groupID)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf(errExec, err)
	}

	if err == sql.ErrNoRows {
		sqlStmt, err = tx.Prepare(requestInsertGroup())
		if err != nil {
			return 0, fmt.Errorf(errStmt, err)
		}
		defer func() { _ = sqlStmt.Close() }()

		err = tx.QueryRow(song.Group).Scan(&groupID)
		if err != nil {
			return 0, fmt.Errorf(errStmt, err)
		}
	}

	sqlStmt, err = tx.Prepare(requestInsertSong())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	startInsert := time.Now()
	err = sqlStmt.QueryRow(
		groupID,
		song.Song,
		song.ReleaseDateTime,
		song.Lyrics,
		song.Link,
	).Scan(&songID)
	endInsert := time.Now()

	if err != nil {
		return 0, fmt.Errorf(errExec, err)
	}

	log.Printf(msgTimeInsert, requestID, endInsert.Sub(startInsert))
	return songID, nil
}
