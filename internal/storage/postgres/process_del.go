package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func (db *DataBase) Delete(id int, requestID string) (err error) {
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

	err = removeSong(tx, id, requestID)
	if err != nil {
		return err
	}

	return nil
}

func removeSong(tx *sql.Tx, id int, requestID string) error {
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

	return nil
}
