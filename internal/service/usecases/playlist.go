package usecases

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
)

func GetPlaylists(input models.GetPlaylistsInput) ([]models.Playlist, error) {
	if check := utility.IsValidUUID(input.UserId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	if check := input.Offset >= 0; !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"offset\""}
	}
	playlists, err :=
		abstractions.PlaylistRepositoryInstance.GetPlaylists(input.UserId, input.Offset)
	if err != nil {
		return nil, err
	}
	return playlists, nil
}

func GetPlaylistSongs(input models.GetPlaylistSongsInput) ([]models.Song, error) {
	if check := utility.IsValidUUID(input.UserId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid user \"uuid\""}
	}
	if check := utility.IsValidUUID(input.PlaylistId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid chat \"uuid\""}
	}

	ownerId, err := abstractions.PlaylistRepositoryInstance.
		GetUserIdByPlaylist(input.PlaylistId)
	if err != nil {
		return nil, err
	}
	if ownerId != input.UserId {
		return nil, &customerrors.ErrInvalidInput{
			Message: "you are not the owner"}
	}

	if input.PlaylistSongId == 0 {
		input.PlaylistSongId = 9223372036854775807
	}

	songs, err :=
		abstractions.PlaylistSongRepositoryInstance.
			GetPlaylistSongs(input.PlaylistId, input.PlaylistSongId)
	if err != nil {
		return nil, err
	}

	return songs, nil

}

func DeletePlaylist(input models.DeletePlaylistInput) error {
	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid user \"id\""}
	}
	if check := utility.IsValidUUID(input.PlaylistId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid playlist\"id\""}
	}

	ownerId, err := abstractions.PlaylistRepositoryInstance.
		GetUserIdByPlaylist(input.PlaylistId)
	if err != nil {
		return err
	}

	if ownerId != input.UserId {
		return &customerrors.ErrInvalidInput{
			Message: "you are not the owner"}
	}

	err = abstractions.PlaylistRepositoryInstance.
		DeletePlaylist(input.PlaylistId)

	return err
}

func AddPlaylist(input models.AddPlaylistInput) error {
	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"id\""}
	}

	valid := utility.IsValidStructField(input, "Name")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"name\" (must be 1-100 chars long)"}
	}

	err := abstractions.PlaylistRepositoryInstance.
		AddPlaylist(models.Playlist{
			Id:   utility.GenerateUUID(),
			Name: input.Name,
		})

	return err
}

func AddPlaylistSong(input models.AddPlaylistSongInput) error {
	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid user \"id\""}
	}
	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid song \"id\""}
	}
	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid playlist \"id\""}
	}

	ownerId, err := abstractions.PlaylistRepositoryInstance.
		GetUserIdByPlaylist(input.PlaylistId)
	if err != nil {
		return err
	}
	if ownerId != input.UserId {
		return &customerrors.ErrInvalidInput{
			Message: "you are not the owner"}
	}

	check, err := abstractions.PlaylistSongRepositoryInstance.
		IsSongInPlaylist(input.PlaylistId, input.SongId)
	if err != nil {
		return err
	}
	if check {
		return &customerrors.ErrInvalidInput{
			Message: "the song is already in a playlist"}
	}

	err = abstractions.PlaylistSongRepositoryInstance.
		AddPlaylistSong(input.PlaylistId, input.SongId)

	return err
}

func DeletePlaylistSong(input models.DeletePlaylistSongInput) error {
	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid user \"id\""}
	}
	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid song \"id\""}
	}
	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid playlist \"id\""}
	}

	ownerId, err := abstractions.PlaylistRepositoryInstance.
		GetUserIdByPlaylist(input.PlaylistId)
	if err != nil {
		return err
	}
	if ownerId != input.UserId {
		return &customerrors.ErrInvalidInput{
			Message: "you are not the owner"}
	}

	check, err := abstractions.PlaylistSongRepositoryInstance.
		IsSongInPlaylist(input.PlaylistId, input.SongId)
	if err != nil {
		return err
	}
	if check {
		return &customerrors.ErrInvalidInput{
			Message: "the song is not in a playlist"}
	}

	err = abstractions.PlaylistSongRepositoryInstance.
		DeletePlaylistSong(input.PlaylistId, input.SongId)

	return err
}
