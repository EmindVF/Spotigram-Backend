package abstractions

import (
	"spotigram/internal/service/models"
)

type UserRepository interface {
	// Adds user to the repository.
	// May return ErrInternal or ErrInvalidInput on failure.
	AddUser(user models.User) error

	// Returns uuid and hashed password of an user by its email.
	// Email validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetUUIDAndPasswordByEmail(email string) (string, string, error)

	// Returns a user's *hashed* password by its uuid.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetPassword(uuid string) (string, error)

	// Returns a user's public key by its uuid.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetPublicKey(uuid string) (string, error)

	// Returns a user's picture by its uuid.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetPicture(uuid string) ([]byte, error)

	// Returns a user by its uuid.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetUser(uuid string) (*models.User, error)

	// Returns a user list.
	// Offset validation is provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetUsers(offset int) ([]models.User, error)

	// Updates a user's name.
	// UUID and name validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	UpdateName(uuid string, name string) error

	// Updates a user's password.
	// UUID and password validation is not provided.
	// *Password is meant to be hashed*.
	// May return ErrInternal or ErrNotFound on failure.
	UpdatePassword(uuid string, password string) error

	// Updates a user's public key.
	// UUID and public key validation is not provided.
	// *Public key is meant to be base64 encoded*.
	// May return ErrInternal or ErrNotFound on failure.
	UpdatePublicKey(uuid string, public_key string) error

	// Updates a user's picture.
	// UUID and image file validation is not provided.
	// *Image is meant to be webp of 512 by 512 or lower size*.
	// May return ErrInternal or ErrNotFound on failure.
	UpdatePicture(uuid string, image []byte) error

	// Returns bool on whether the user uuid is present.
	// UUID validation is not provided.
	// May return ErrInternal on failure.
	DoesUserExist(uuid string) (bool, error)
}

type FriendRepository interface {
	// Returns a user's friends by its uuid.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetFriends(uuid string, offset int) ([]models.Friend, error)

	// Returns a char id with between two user.
	// UUID sort is provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetChatIdByFriend(uuid1, uuid2 string) (string, error)

	// Returns a friend by chat id.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetFriendByChatId(uuid string) (f *models.Friend, e error)

	// Adds friend to the repository.
	// UUID sort is provided.
	// May return ErrInternal or ErrInvalidInput on failure.
	AddFriend(models.Friend) error

	// Deletes a friend from the repository.
	// UUID sort is provided.
	// May return ErrInternal, ErrNotFound or ErrInvalidInput on failure.
	DeleteFriend(uuid1, uuid2 string) error

	// Returns bool on whether the friend is present.
	// UUID validation is not provided.
	// UUID sort is provided.
	// May return ErrInternal on failure.
	DoesFriendExist(uuid1, uuid2 string) (bool, error)
}

type FriendRequestRepository interface {
	// Adds friend request to the repository.
	// May return ErrInternal or ErrInvalidInput on failure.
	AddFriendRequest(user models.FriendRequest) error

	// Update friend request in the repository.
	// May return ErrInternal or ErrNotFound on failure.
	UpdateIsIgnored(senderUUID, recipientUUID string, isIgnored bool) error

	// Deletes friend request from the repository.
	// May return ErrInternal, ErrInvalidInput or ErrNotFound on failure.
	DeleteFriendRequest(senderUUID, recipientUUID string) error

	// Returns a user's sent friend requests by its uuid.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetFriendRequestsSent(uuid string, offset int) ([]models.FriendRequest, error)

	// Returns a user's received friend requests by its uuid.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetFriendRequestsReceived(uuid string, offset int) ([]models.FriendRequest, error)

	// Returns bool on whether the friend is present.
	// UUID validation is not provided.
	// May return ErrInternal on failure.
	DoesFriendRequestExist(senderUUID, recipientUUID string) (bool, error)
}

type ChatRepository interface {
	// Deletes a whole chat from the repository by its id.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	DeleteChat(uuid string) error

	// Returns first 100 messages before a certain time id.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetMessages(chatId string, timeId int64) ([]models.Message, error)

	// Adds a message to the repository.
	// UUID validation is not provided.
	// May return ErrInternal on failure.
	AddMessage(message models.Message) error

	// Deletes a message from the chat by chat and message time ids.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	DeleteMessage(chatId string, timeId int64) error
}

type PlaylistRepository interface {
	// Returns first 100 messages before a certain time id.
	// UUID validation is not provided.
	// May return ErrInternal.
	GetPlaylists(userId string, offset int) ([]models.Playlist, error)

	// Returns the playlist by its id.
	// UUID validation is not provided.
	GetPlaylist(playlistId string) (*models.Playlist, error)

	// Returns the playlist by its id.
	// UUID validation is not provided
	UpdatePlaylist(models.Playlist) error

	// Adds a playlist.
	// May return ErrInternal or ErrInvalidInput.
	AddPlaylist(models.Playlist) error

	// Deletes a playlist.
	// May return ErrInternal or ErrNotFound.
	DeletePlaylist(playlistId string) error
}

type PlaylistSongRepository interface {
	// Checks if a song already is in a playlist.
	// May return ErrInternal of ErrNotFound
	IsSongInPlaylist(playlistId string, songId string) (bool, error)

	// Adds a song to a playlist.
	// May return ErrInternal.
	AddPlaylistSong(playlistSong models.PlaylistSong) error

	// Deletes a song from a playlist.
	// May return ErrInternal.
	DeletePlaylistSong(playlistId string, songId string) error

	// Adds all songs to a playlist.
	// May return ErrInternal.
	DeletePlaylistSongs(playlistId string) error

	// Returns first 100 playlist songs.
	// UUID validation is not provided
	// May return Err internal
	GetPlaylistSongs(playlistId string) ([]models.Song, error)
}

type SongRepository interface {
	// Returns first 100 messages before a certain time id.
	// UUID validation is not provided
	// May return ErrInternal
	GetSongs(offset int) ([]models.Song, error)

	// Returns songs info by its id.
	// UUID validation is not provided.
	// May return ErrInternal.
	GetSongInfo(songId string) (*models.Song, error)

	// Returns songs file by its id.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetSongFile(songId string) ([]byte, error)

	// Returns a song's picture by its uuid.
	// UUID validation is not provided.
	// May return ErrInternal or ErrNotFound on failure.
	GetPicture(songId string) ([]byte, error)

	// Adds song to the repository.
	// May return ErrInternal or ErrInvalidInput on failure.
	AddSong(song models.Song, picture []byte, file []byte) error

	// Deletes a song from the repository.
	// May return ErrInternal or ErrNotFound on failure.
	DeleteSong(songId string) error
}

type SongChunkRepository interface {
	// Uploads song chunk to the repository.
	// May return ErrInternal on failure.
	AddSongChunk(songId string, id int, chunk []byte) error

	// Uploads a song chunks to the repository.
	// May return ErrInternal or ErrNotFound on failure.
	GetSongChunk(songId string, index int) ([]byte, error)

	// Deletes a song from the repository.
	// May return ErrInternal or ErrNotFound on failure.
	DeleteSongChunks(songId string) error
}
