package models

import "encoding/json"

// Fields not marked with tags are meant to be set manually.
// Fields markes with 'json' tag are meant to be unmarshalled from the request.

// Auth
type SignUpInput struct {
	Name              string `validate:"required,min=5,max=100" json:"name"`
	Email             string `validate:"required,max=100,email" json:"email"`
	Password          string `validate:"required,min=8,max=72" json:"password"`
	PasswordConfirmed string `validate:"required,min=8,max=72" json:"password_confirmed"`
}

type SignInInput struct {
	Email    string `validate:"required,max=100,email" json:"email"`
	Password string `validate:"required,min=8,max=72" json:"password"`
}

type LogoutInput struct {
	AccessTokenUUID string
	RefreshToken    string
}

type RefreshAccessTokenInput struct {
	RefreshToken string
}

// Deserialization
type DeserializeTokenInput struct {
	AccessToken string
}

// User
type GetUsersInput struct {
	Offset         int    `validate:"required" json:"offset"`
	UserNameFilter string `validate:"max=100" json:"username_filter"`
}

type GetUserInfoInput struct {
	UserUUID string `validate:"required,min=8,max=130" json:"id"`
}

type GetPublicKeyInput struct {
	UserUUID string `validate:"required,min=8,max=130" json:"id"`
}

type GetPictureInput struct {
	UserUUID string `validate:"required,min=8,max=130" json:"id"`
}

// Me
type ChangeNameInput struct {
	Name     string `validate:"required,min=5,max=100" json:"name"`
	UserUUID string
}

type ChangePasswordInput struct {
	OldPassword          string `validate:"required,min=8,max=72" json:"old_password"`
	NewPassword          string `validate:"required,min=8,max=72" json:"new_password"`
	NewPasswordConfirmed string `validate:"required,min=8,max=72" json:"new_password_confirmed"`
	UserUUID             string
}

type ChangePublicKeyInput struct {
	PublicKey string `validate:"required,min=1,max=6120" json:"public_key"`
	UserUUID  string
}

type ChangePictureInput struct {
	Image    []byte
	UserUUID string
}

// Friends
type GetFriendsInput struct {
	UserUUID string
	Offset   int `validate:"required" json:"offset"`
}

type DeleteFriendInput struct {
	User1UUID string
	User2UUID string `validate:"required,min=8,max=130" json:"id"`
}

// Friend requests
type GetFriendRequestsSentInput struct {
	UserUUID string
	Offset   int `validate:"required" json:"offset"`
}

type GetFriendRequestsReceivedInput struct {
	UserUUID string
	Offset   int `validate:"required" json:"offset"`
}

type AddFriendRequestInput struct {
	SenderUUID    string
	RecipientUUID string `validate:"required,min=8,max=130" json:"id"`
}

type UpdateFriendRequestInput struct {
	SenderUUID    string `validate:"required,min=8,max=130" json:"id"`
	RecipientUUID string
	IsIgnored     bool `json:"is_ignored"`
}

type DeleteFriendRequestInput struct {
	SenderUUID    string
	RecipientUUID string `validate:"required,min=8,max=130" json:"id"`
}

type AcceptFriendRequestInput struct {
	SenderUUID    string `validate:"required,min=8,max=130" json:"id"`
	RecipientUUID string
}

// Websocket
type WebsocketPayload struct {
	Action  string          `json:"action"`
	Content json.RawMessage `json:"content"`
}

// Chat
type GetMessagesInput struct {
	UserId string `json:"-"`
	ChatId string `json:"chat_id"`
	TimeId int64  `json:"id"`
}

type GetUnreadMessagesInput struct {
	UserId string `json:"-"`
	ChatId string `json:"chat_id"`
}

type SendMessageInput struct {
	UserId      string `json:"-"`
	ChatId      string `json:"chat_id"`
	Content     string `json:"content"`
	TimeId      int64  `json:"id"`
	IsEncrypted bool   `json:"is_encrypted"`
}

type DeleteMessageInput struct {
	UserId string `json:"-"`
	ChatId string `json:"chat_id"`
	TimeId int64  `json:"id"`
}

type CheckStatusInput struct {
	UserId string `json:"user_id"`
}

type UpdateStatusInput struct {
	UserId string `json:"user_id"`
	Status int    `json:"status"`
}

// Songs
type GetSongsInput struct {
	Offset          int    `validate:"required" json:"offset"`
	SongNameFilter  string `validate:"max=100" json:"songname_filter"`
	CreatorIdFilter string `validate:"max=130" json:"creatorid_filter"`
}

type GetSongInfoInput struct {
	SongId string `json:"id"`
}

type GetSongFileInput struct {
	SongId string `json:"id"`
}

type GetSongPictureInput struct {
	SongId string `json:"id"`
}

type GetSongChunkInput struct {
	FileName string
}

type AddSongInput struct {
	UserId string
	Name   string `validate:"required,min=5,max=100" json:"name"`
	File   []byte
}

type UpdateSongNameInput struct {
	UserId string
	SongId string `json:"id"`
	Name   string `validate:"required,min=5,max=100" json:"new_name"`
}

type DeleteSongInput struct {
	UserId string
	SongId string `json:"id"`
}

// Playlists
type GetPlaylistsInput struct {
	UserId             string
	Offset             int    `validate:"required" json:"offset"`
	PlaylistNameFilter string `validate:"max=100" json:"playlistname_filter"`
}

type AddPlaylistInput struct {
	UserId string
	Name   string `validate:"required,min=5,max=100" json:"name"`
}

type UpdatePlaylistNameInput struct {
	UserId     string
	PlaylistId string `json:"id"`
	Name       string `validate:"required,min=5,max=100" json:"new_name"`
}

type DeletePlaylistInput struct {
	UserId     string
	PlaylistId string `json:"id"`
}

type AddPlaylistSongInput struct {
	UserId     string
	PlaylistId string `json:"id"`
	SongId     string `json:"song_id"`
}

type DeletePlaylistSongInput struct {
	UserId     string
	PlaylistId string `json:"id"`
	SongId     string `json:"song_id"`
}

type GetPlaylistSongsInput struct {
	UserId     string
	PlaylistId string `json:"id"`
}

// Read Times
type GetReadTimeInput struct {
	UserId string `json:"user_id"`
	ChatId string `json:"chat_id"`
}

type UpdateReadTimeInput struct {
	UserId string `json:"-"`
	ChatId string `json:"chat_id"`
	TimeId int64  `json:"time_id"`
}
