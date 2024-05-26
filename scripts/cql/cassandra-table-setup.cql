CREATE TABLE IF NOT EXISTS messages (
    time_id bigint,
    creation_date timestamp,
    user_id UUID,
    chat_id UUID,
    content text,
    encrypted boolean,
    PRIMARY KEY (chat_id, time_id)
) WITH CLUSTERING ORDER BY (time_id DESC);

CREATE TABLE IF NOT EXISTS playlist_songs (
    playlist_id UUID,
    song_id UUID,
    PRIMARY KEY (playlist_id, song_id)
);

CREATE TABLE IF NOT EXISTS song_chunks (
    song_id UUID,
    ind int,
    file blob,
    PRIMARY KEY (song_id, ind)
) WITH CLUSTERING ORDER BY (ind ASC);