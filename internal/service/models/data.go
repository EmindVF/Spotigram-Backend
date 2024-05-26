package models

import "time"

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Verified bool   `json:"verified"`
}

type Friend struct {
	Id1    string `json:"user_id1"`
	Id2    string `json:"user_id2"`
	ChatId string `json:"chat_id"`
}

type FriendRequest struct {
	SenderId    string `json:"sender_id"`
	RecipientId string `json:"recipient_id"`
	IsIgnored   bool   `json:"is_ignored"`
}

type Message struct {
	UserId      string    `json:"user_id"`
	ChatId      string    `json:"chat_id"`
	Content     string    `json:"content"`
	Date        time.Time `json:"date"`
	TimeId      int64     `json:"id"`
	IsEncrypted bool      `json:"is_encrypted"`
}

type Playlist struct {
	Id     string `json:"id"`
	UserId string `json:"-"`
	Name   string `json:"name"`
}

type PlaylistSong struct {
	PlaylistId string
	SongId     string `json:"song_id"`
}

type Song struct {
	Id        string `json:"id"`
	CreatorId string `json:"creator_id"`
	Name      string `json:"name"`
	Length    int    `json:"length"`
	Streams   int    `json:"streams"`
}

type LastStream struct {
	UserId string
	SongId string
	Time   int64
}

type ReadTime struct {
	UserId string
	ChatId string
	TimeId int64
}
