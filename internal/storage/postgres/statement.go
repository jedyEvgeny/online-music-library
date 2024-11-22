package storage

func requestSelectGroupID() string {
	return `
	SELECT g_id 
	FROM music_groups 
	WHERE name=$1
	`
}

func requestInsertGroup() string {
	return `
	INSERT INTO music_groups (name) 
	VALUES ($1) 
	RETURNING g_id
	`
}

func selectSongByGroup() string {
	return `
	SELECT s_id 
	FROM songs 
	WHERE group_id = $1 AND song = $2`
}

func reqUpdateSongToinitState() string {
	return `
    UPDATE songs
    SET release_date = $1, 
        lyrics = $2, 
        link = $3, 
        updated_at = NOW()
    WHERE s_id = $4`
}

func requestInsertSong() string {
	return `
	INSERT INTO songs 
	(group_id, song, release_date, lyrics, link)
    VALUES ($1, $2, $3, $4, $5)
	RETURNING s_id
	`
}

func requestDeleteSong() string {
	return `
	DELETE FROM songs 
	WHERE s_id = $1
	`
}
