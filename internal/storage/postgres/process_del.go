package storage

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func (db *DataBase) Delete(id int, requestID string) (err error) {
	const op = "Delete"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))
	startTx := time.Now()

	tx, err := db.db.Begin()
	if err != nil {
		return fmt.Errorf(errCreateTx, err)
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

	err = db.removeSong(tx, id, requestID)
	if err != nil {
		return err
	}
	endTx := time.Now()

	db.log.Debug(fmt.Sprintf(logTxTime, requestID, endTx.Sub(startTx)))
	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return nil
}

func (db *DataBase) removeSong(tx *sql.Tx, id int, requestID string) error {
	const op = "removeSong"
	db.log.Debug(fmt.Sprintf(logStart, requestID, op))

	sqlStmt, err := tx.Prepare(requestDeleteSong())
	if err != nil {
		return fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	result, err := sqlStmt.Exec(id)
	if err != nil {
		return fmt.Errorf(errExec, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf(errRes, err)
	}

	if rowsAffected == 0 {
		log.Printf(msgResAffected, requestID, id)
	}

	db.log.Debug(fmt.Sprintf(logEnd, requestID, op))
	return nil
}
