CREATE TABLE IF NOT EXISTS messages (
    time_id bigint,
    user_id UUID,
    chat_id UUID,
    content text,
    encrypted boolean,
    PRIMARY KEY (chat_id, time_id)
) WITH CLUSTERING ORDER BY (time_id DESC);

CREATE TABLE IF NOT EXISTS playlist_songs (
    time_id bigint,
    playlist_id UUID,
    song_id UUID,
    PRIMARY KEY (playlist_id, song_id)
);