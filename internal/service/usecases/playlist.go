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

	playlist, err := abstractions.PlaylistRepositoryInstance.
		GetPlaylist(input.PlaylistId)
	if err != nil {
		return nil, err
	}
	if playlist.UserId != input.UserId {
		return nil, &customerrors.ErrInvalidInput{
			Message: "you are not the owner"}
	}

	songs, err :=
		abstractions.PlaylistSongRepositoryInstance.
			GetPlaylistSongs(input.PlaylistId)
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

	playlist, err := abstractions.PlaylistRepositoryInstance.
		GetPlaylist(input.PlaylistId)
	if err != nil {
		return err
	}

	if playlist.UserId != input.UserId {
		return &customerrors.ErrInvalidInput{
			Message: "you are not the owner"}
	}

	err = abstractions.PlaylistRepositoryInstance.
		DeletePlaylist(input.PlaylistId)
	if err != nil {
		return err
	}

	err = abstractions.PlaylistSongRepositoryInstance.
		DeletePlaylistSongs(input.PlaylistId)

	return err
}

func AddPlaylist(input models.AddPlaylistInput) (string, error) {
	if check := utility.IsValidUUID(input.UserId); !check {
		return "", &customerrors.ErrInvalidInput{
			Message: "invalid \"id\""}
	}

	valid := utility.IsValidStructField(input, "Name")
	if !valid {
		return "", &customerrors.ErrInvalidInput{
			Message: "invalid \"name\" (must be 5-100 chars long)"}
	}

	uuid := utility.GenerateUUID()

	err := abstractions.PlaylistRepositoryInstance.
		AddPlaylist(models.Playlist{
			Id:   uuid,
			Name: input.Name,
		})
	if err != nil {
		return "", err
	}
	return uuid, nil
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

	playlist, err := abstractions.PlaylistRepositoryInstance.
		GetPlaylist(input.PlaylistId)
	if err != nil {
		return err
	}
	if playlist.UserId != input.UserId {
		return &customerrors.ErrInvalidInput{
			Message: "you are not the owner"}
	}

	if playlist.Length == 100 {
		return &customerrors.ErrInvalidInput{
			Message: "the playlist is at max length"}
	}

	song, err := abstractions.SongRepositoryInstance.
		GetSongInfo(input.SongId)
	if err != nil {
		return err
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

	playlist.Length++

	err = abstractions.PlaylistRepositoryInstance.
		UpdatePlaylist(*playlist)
	if err != nil {
		return err
	}

	err = abstractions.PlaylistSongRepositoryInstance.
		AddPlaylistSong(models.PlaylistSong{
			PlaylistId: input.PlaylistId,
			SongId:     song.Id,
			UserId:     song.CreatorId,
			Name:       song.Name,
			Length:     song.Length,
		})

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

	playlist, err := abstractions.PlaylistRepositoryInstance.
		GetPlaylist(input.PlaylistId)
	if err != nil {
		return err
	}
	if playlist.UserId != input.UserId {
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

	playlist.Length--

	err = abstractions.PlaylistRepositoryInstance.
		UpdatePlaylist(*playlist)
	if err != nil {
		return err
	}

	err = abstractions.PlaylistSongRepositoryInstance.
		DeletePlaylistSong(input.PlaylistId, input.SongId)

	return err
}
