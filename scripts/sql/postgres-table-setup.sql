CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id UUID NOT NULL PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  email VARCHAR(100) NOT NULL UNIQUE,
  password BYTEA NOT NULL,
  picture BYTEA,
  verified BOOLEAN,
  public_key BYTEA
);

CREATE TABLE IF NOT EXISTS friend_requests (
  sender_id UUID NOT NULL,
  recipient_id UUID NOT NULL,
  is_ignored BOOLEAN,
  PRIMARY KEY (sender_id, recipient_id),
  FOREIGN KEY (sender_id) REFERENCES users (id),
  FOREIGN KEY (recipient_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS friendships (
  user1_id UUID NOT NULL,
  user2_id UUID NOT NULL,
  chat_id UUID NOT NULL UNIQUE,
  PRIMARY KEY (user1_id, user2_id),
  FOREIGN KEY (user1_id) REFERENCES users (id),
  FOREIGN KEY (user2_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS songs (
  id UUID NOT NULL PRIMARY KEY,
  creator_id UUID REFERENCES users (id),
  name VARCHAR(100) NOT NULL,
  length INTEGER NOT NULL,
  streams INTEGER,
  picture BYTEA,
  file BYTEA
);

CREATE TABLE IF NOT EXISTS playlists (
  id UUID NOT NULL PRIMARY KEY,
  name VARCHAR(100),
  user_id UUID REFERENCES users (id) 
);

CREATE TABLE IF NOT EXISTS read_times (
  user_id UUID NOT NULL,
  chat_id UUID NOT NULL,
  time_id BIGINT,
  PRIMARY KEY (user_id, chat_id),
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (chat_id) REFERENCES friendships (chat_id)
);