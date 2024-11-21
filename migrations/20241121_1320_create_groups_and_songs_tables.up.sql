CREATE TABLE music_groups (
    g_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_music_groups_name ON music_groups (name);

CREATE TABLE songs (
    s_id SERIAL PRIMARY KEY,
    group_id INT NOT NULL REFERENCES music_groups(g_id) ON DELETE CASCADE,
    song TEXT NOT NULL,
    release_date DATE NOT NULL,
    lyrics TEXT NOT NULL,
    link TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_songs_group_id ON songs (group_id);
CREATE INDEX idx_songs_song ON songs (song);
CREATE INDEX idx_songs_release_date ON songs (release_date);
