package usecases

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
)

// A use case to get the songs list.
func GetSongs(gsi models.GetSongsInput) ([]models.Song, error) {
	if check := gsi.Offset >= 0; !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"offset\""}
	}
	songs, err :=
		abstractions.SongRepositoryInstance.GetSongs(gsi.Offset)
	if err != nil {
		return nil, err
	}
	return songs, nil
}

// A use case to delete a song.
// Expects access token deserialization beforehand.
// Validates the passed uuid.
// May return ErrInvalidInput, ErrInternal, ErrNotFound on failure.
func DeleteSong(input models.DeleteSongInput) error {
	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	if check := utility.IsValidUUID(input.SongId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}

	err := abstractions.SongRepositoryInstance.
		DeleteSong(input.SongId)
	return err
}

// A use case to upload a song
// Returns the uuid of the recipient
func AddSong(input models.AddSongInput) error {
	/*
		if check := utility.IsValidUUID(input.UserId); !check {
			return "", &customerrors.ErrInvalidInput{
				Message: "invalid user \"uuid\""}
		}

		if len(input.Content) > 2048 || len(input.Content) == 0 {
			return "", &customerrors.ErrInvalidInput{
				Message: "message is too long"}
		}
		friend, err := abstractions.FriendRepositoryInstance.GetFriendByChatId(input.ChatId)
		if err != nil {
			return "", err
		}

		var uuidRecipient string

		if input.UserId == friend.Id1 {
			uuidRecipient = friend.Id2
			input.TimeId = time.Now().UnixMicro()*10 + 1
		} else if input.UserId == friend.Id2 {
			uuidRecipient = friend.Id1
			input.TimeId = time.Now().UnixMicro()*10 + 2
		} else {
			return "", &customerrors.ErrInvalidInput{
				Message: "user is not a member of this chat"}
		}

		err = abstractions.ChatRepositoryInstance.AddMessage(input)
		if err != nil {
			return "", nil
		}

		return uuidRecipient, nil*/
	return nil
}
