CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id UUID NOT NULL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  password VARCHAR(100) NOT NULL,
  picture BYTEA,
  verified BOOLEAN,
  public_key VARCHAR(4096)
);

CREATE TABLE IF NOT EXISTS chats (
  id UUID NOT NULL PRIMARY KEY,
  message_count integer 
);

CREATE TABLE IF NOT EXISTS friendships (
  user1_id UUID NOT NULL,
  user2_id UUID NOT NULL,
  chat_id UUID NOT NULL UNIQUE,
  PRIMARY KEY (user1_id,user2_id),
  FOREIGN KEY (user1_id) REFERENCES users (id),
  FOREIGN KEY (user2_id) REFERENCES users (id),
  FOREIGN KEY (chat_id) REFERENCES chats (id)
);

CREATE TABLE IF NOT EXISTS songs (
  id UUID NOT NULL PRIMARY KEY,
  creator_id UUID REFERENCES users (id),
  name VARCHAR(100) NOT NULL,
  path VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS playlists (
  id UUID NOT NULL PRIMARY KEY,
  name VARCHAR(100),
  user_id UUID REFERENCES users (id) 
);

CREATE TABLE IF NOT EXISTS playlist_songs (
  playlist_id UUID,
  song_id UUID,
  PRIMARY KEY (playlist_id,song_id),
  FOREIGN KEY (playlist_id) REFERENCES playlists (id),
  FOREIGN KEY (song_id) REFERENCES songs (id)
);

CREATE TABLE IF NOT EXISTS messages (
  id UUID NOT NULL PRIMARY KEY,
  chat_id UUID,
  user_id UUID,
  text VARCHAR(500), 
  sending_date TIMESTAMP,
  FOREIGN KEY (chat_id) REFERENCES chats (id),
  FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS friend_requests (
  from_id UUID NOT NULL,
  to_id UUID NOT NULL,
  PRIMARY KEY (from_id, to_id),
  FOREIGN KEY (from_id) REFERENCES users (id),
  FOREIGN KEY (to_id) REFERENCES users (id)
)
