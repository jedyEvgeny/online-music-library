package storage

func reqSelectGroupID() string {
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

func requestInsertSong() string {
	return `
	INSERT INTO songs 
	(group_id, song, release_date, lyrics, link)
    VALUES ($1, $2, $3, $4, $5)
	RETURNING s_id
	`
}
