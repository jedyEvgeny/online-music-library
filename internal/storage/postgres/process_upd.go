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
	fmt.Println("выполняем checkSong")

	currGroupID, err := db.findGroupID(tx, songID)
	if err != nil {
		return err
	}
	fmt.Println("выполнили currGroup без ошибки")

	groupIdLikeGroupInNewSong, err := db.findExistGroupIdLikeNewGroup(tx, e.Group, requestID)
	if err != nil {
		return fmt.Errorf(errNoContent, err)
	}
	fmt.Println("выполнили groupIdLikeGroupInNewSong без ошибки")

	foundMoreThanOneSongForCurrGroup, err := db.hasMoreThanOneSongForCurrGroup(tx, currGroupID, requestID)
	if err != nil {
		return fmt.Errorf(errNoContent, err)
	}

	if e.Group != "" {
		err = updGroup(tx, e, foundMoreThanOneSongForCurrGroup, currGroupID, songID, groupIdLikeGroupInNewSong)
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

func findSongIdBySongID(tx *sql.Tx, songID int) error {
	sqlStmt, err := tx.Prepare(requestSelectSongID())
	if err != nil {
		return fmt.Errorf(errStmt, err)
	}
	defer func() { _ = sqlStmt.Close() }()

	var countID int
	err = sqlStmt.QueryRow(songID).Scan(&countID)
	if err != nil {
		return fmt.Errorf(errExec, err)
	}
	if countID == 0 {
		return fmt.Errorf(errNoContent)
	}
	return nil
}

// findGroupID возвращает id группы из таблицы songs
// и косвенно проверяет, существует ли переданный songID
func (db *DataBase) findGroupID(tx *sql.Tx, songID int) (int, error) {
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

	return currGroupID, nil
}

func (db *DataBase) findExistGroupIdLikeNewGroup(tx *sql.Tx, group, requestID string) (int, error) {
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
	return groupIdLikeGroupInNewSong, nil
}

func (db *DataBase) hasMoreThanOneSongForCurrGroup(tx *sql.Tx, groupID int, requestID string) (bool, error) {
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

func updGroup(tx *sql.Tx, e *service.EnrichedSong, foundMoreTwoSongsByCurrGroup bool, groupID, songID, groupIdLikeGroupInNewSong int) error {
	//fmt.Println("Выполняем updGround")

	//Случай, когда новой группы нет в таблице music_groups, но есть больше 1 песни у старой группы или
	//Случай, когда у существующей группы больше 2х песен и новой группы нет в БД
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

	//Случай, когда у существующей группы одна песня, но новая группа уже есть в БД или
	//Случай, когда у существующей группы больше 2х песен и новая группа есть в БД
	if (!foundMoreTwoSongsByCurrGroup && groupIdLikeGroupInNewSong != 0) ||
		(foundMoreTwoSongsByCurrGroup && groupIdLikeGroupInNewSong != 0) {
		err := updGroupIdInSongsBySongID(tx, groupIdLikeGroupInNewSong, songID)
		if err != nil {
			return err
		}
	}

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
