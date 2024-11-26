package storage

import (
	"bytes"
	"database/sql"
	"fmt"
	"jedyEvgeny/online-music-library/internal/app/service"
	"net/http"
)

func (db *DataBase) Update(song *service.EnrichedSong, songID int, requestID string) (statusCode int, err error) {
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

	return http.StatusOK, nil
}

func (db *DataBase) checkSong(tx *sql.Tx, e *service.EnrichedSong, songID int, requestID string) error {
	groupIDInMusicGroupsTable, err := db.findGroupID(tx, e.Group)
	if err != nil {
		return fmt.Errorf(errNoContent, err)
	}

	groupIDInSongTable, err := db.isExistingSong(tx, songID, requestID)
	if err != nil {
		return fmt.Errorf(errNoContent, err)
	}

	foundMoreTwoSongsByCurrGroup, err := db.isMoreThenTwoSongsCurrGroup(tx, groupIDInSongTable, requestID)
	if err != nil {
		return fmt.Errorf(errNoContent, err)
	}

	if e.Group != "" {
		err = updGroup(tx, e, foundMoreTwoSongsByCurrGroup, groupIDInMusicGroupsTable, groupIDInSongTable, songID)
		if err != nil {
			return err
		}
	}

	needToChange, err := updSong(tx, e, songID)
	if !needToChange && err == nil { //если поменялось только название исполнителя
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}

func (db *DataBase) findGroupID(tx *sql.Tx, group string) (int, error) {
	var groupIDInMusicGroupsTable int
	sqlStmt, err := tx.Prepare(requestSelectGroupID())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(group).Scan(&groupIDInMusicGroupsTable)
	if err != nil && err != sql.ErrNoRows {
		return 0, fmt.Errorf(errExec, err)
	}

	return groupIDInMusicGroupsTable, nil
}

func (db *DataBase) isExistingSong(tx *sql.Tx, songID int, requestID string) (int, error) {
	var groupIDInSongTable int
	sqlStmt, err := tx.Prepare(requestSelectGroupBySong())
	if err != nil {
		return 0, fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(songID).Scan(&groupIDInSongTable)
	if err != nil {
		return 0, fmt.Errorf(errExec, err)
	}
	return groupIDInSongTable, nil
}

func (db *DataBase) isMoreThenTwoSongsCurrGroup(tx *sql.Tx, groupID int, requestID string) (bool, error) {
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
	return true, nil
}

func updGroup(tx *sql.Tx, e *service.EnrichedSong, foundMoreTwoSongsByCurrGroup bool, groupIDInMusicGroupsTable, groupIDInSongTable, songID int) error {
	//Если ранее не было нового исполнителя имеющейся композиции,
	//Но есть другие композиции имеющегося исполнителя
	if groupIDInMusicGroupsTable == 0 && foundMoreTwoSongsByCurrGroup {
		err := createGroup(tx, e, groupIDInMusicGroupsTable)
		if err != nil {
			return err
		}

		err = updGroupIdInSongsBySongID(tx, groupIDInMusicGroupsTable, songID)
		if err != nil {
			return err
		}
	}

	//Если ранее не было нового исполнителя имеющейся композиции,
	//и нет других композиций имеющегося исполнителя
	if groupIDInMusicGroupsTable == 0 && !foundMoreTwoSongsByCurrGroup {
		err := updateGroupForSingleSongBySuchGroup(tx, e, groupIDInMusicGroupsTable)
		if err != nil {
			return err
		}
	}

	return nil
}

func createGroup(tx *sql.Tx, e *service.EnrichedSong, groupIDInMusicGroupsTable int) error {
	sqlStmt, err := tx.Prepare(requestInsertGroup())
	if err != nil {
		return fmt.Errorf(errExec, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	err = sqlStmt.QueryRow(e.Group).Scan(&groupIDInMusicGroupsTable)
	if err != nil {
		return err
	}
	return nil
}

func updGroupIdInSongsBySongID(tx *sql.Tx, groupIDInMusicGroupsTable, songID int) error {
	sqlStmt, err := tx.Prepare(requestUpdateGroupIdInSongsBySongID())
	if err != nil {
		return fmt.Errorf(errExec, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	_, err = sqlStmt.Exec(
		groupIDInMusicGroupsTable,
		songID,
	)
	if err != nil {
		return err
	}
	return nil
}

func updateGroupForSingleSongBySuchGroup(tx *sql.Tx, e *service.EnrichedSong, groupIDInMusicGroupsTable int) error {
	sqlStmt, err := tx.Prepare(requestUpdateGroupForSingleSongBySuchGroup())
	if err != nil {
		return fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	_, err = sqlStmt.Exec(e.Group, groupIDInMusicGroupsTable)
	if err != nil {
		return fmt.Errorf(errExec, err)
	}
	return nil
}

func updSong(tx *sql.Tx, e *service.EnrichedSong, songID int) (bool, error) {
	query, args := createQueryUpdSong(e, songID)

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

	return true, nil
}

func createQueryUpdSong(e *service.EnrichedSong, songID int) (string, []interface{}) {
	var buf bytes.Buffer
	var args []interface{}

	buf.WriteString("UPDATE songs SET ")

	if e.Song != "" {
		args = append(args, e.Song)
		buf.WriteString(fmt.Sprintf("song = $%d", len(args)))
	}
	if e.ReleaseDate != "" {
		args = append(args, e.ReleaseDate)
		buf.WriteString(fmt.Sprintf("release_date = $%d", len(args)))
	}
	if e.Lyrics != "" {
		args = append(args, e.Lyrics)
		buf.WriteString(fmt.Sprintf("lyrics = $%d", len(args)))
	}
	if e.Link != "" {
		args = append(args, e.Link)
		buf.WriteString(fmt.Sprintf("link = $%d", len(args)))
	}

	if len(args) > 0 {
		args = append(args, songID)
		buf.WriteString(fmt.Sprintf(" WHERE s_id = $%d", len(args)))
	}

	return buf.String(), args
}
