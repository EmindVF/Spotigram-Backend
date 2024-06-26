package usecases

import (
	"os"
	"os/exec"
	"spotigram/internal/customerrors"
	"spotigram/internal/service/abstractions"
	"spotigram/internal/service/models"
	"spotigram/internal/utility"
	"strconv"
	"strings"
)

// A use case to get the songs list.
func GetSongs(gsi models.GetSongsInput) ([]models.Song, error) {
	if check := gsi.Offset >= 0; !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"offset\""}
	}
	if !utility.IsValidStructField(gsi, "SongNameFilter") {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"songname_filter\" (empty or under 100 chars)"}
	}
	if gsi.CreatorIdFilter != "" {
		if check := utility.IsValidUUID(gsi.CreatorIdFilter); !check {
			return nil, &customerrors.ErrInvalidInput{
				Message: "invalid \"creatorid_filter\""}
		}
	}

	songs, err :=
		abstractions.SongRepositoryInstance.GetSongs(
			gsi.Offset, gsi.SongNameFilter, gsi.CreatorIdFilter)
	if err != nil {
		return nil, err
	}
	return songs, nil
}

// A use case to get the songs list.
func GetSongInfo(gsi models.GetSongInfoInput) (*models.Song, error) {
	if check := utility.IsValidUUID(gsi.SongId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	song, err :=
		abstractions.SongRepositoryInstance.GetSongInfo(gsi.SongId)
	if err != nil {
		return nil, err
	}
	return song, nil
}

// A use case to get the songs list.
func GetSongFile(gsi models.GetSongFileInput) ([]byte, error) {
	if check := utility.IsValidUUID(gsi.SongId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	file, err :=
		abstractions.SongRepositoryInstance.GetSongFile(gsi.SongId)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func GetSongPicture(input models.GetSongPictureInput) ([]byte, error) {
	if check := utility.IsValidUUID(input.SongId); !check {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid \"uuid\""}
	}
	pic, err :=
		abstractions.SongRepositoryInstance.GetPicture(input.SongId)
	if err != nil {
		return nil, err
	}
	return pic, nil
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

	user, err := abstractions.UserRepositoryInstance.
		GetUser(input.UserId)
	if err != nil {
		return err
	}

	if !user.Verified {
		return &customerrors.ErrUnauthorized{
			Message: "you are not verified",
		}
	}

	/*
		err = abstractions.PlaylistSongRepositoryInstance.
			DeleteSong(input.SongId)
		if err != nil {
			return err
		}*/

	err = abstractions.SongRepositoryInstance.
		DeleteSong(input.SongId)
	if err != nil {
		return err
	}
	err = abstractions.SongChunkRepositoryInstance.
		DeleteSongChunks(input.SongId)
	return err
}

// A use case to upload a song
// Returns the uuid of the recipient
func AddSong(input models.AddSongInput) error {

	valid := utility.IsValidStructField(input, "Name")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"name\" (must be 5-100 chars long)"}
	}

	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid user \"uuid\""}
	}
	var albumCover []byte = nil
	if utility.IsValidMP3(input.File) {
		pic, _, err := utility.GetMP3AlbumCover(input.File)
		if err != nil {
			return &customerrors.ErrInternal{
				Message: "cannot read tags of the song",
			}
		}

		if pic != nil {
			if len(pic) == 0 || len(pic) > 2*1024*1024 {
				return &customerrors.ErrInvalidInput{
					Message: "invalid album cover, must be a png or a jpg (under 2 megabytes)"}
			}

			albumCover, err = utility.ConvertAndResizeImageToWebP(pic, 512, 512)
			if err != nil {
				return &customerrors.ErrInvalidInput{
					Message: "invalid album cover, must be a png or a jpg (under 2 megabytes)"}
			}
		}
	}

	song := models.Song{
		Id:        utility.GenerateUUID(),
		CreatorId: input.UserId,
		Name:      input.Name,
		Length:    0,
	}

	// Split the song into chunks
	tempDir, err := os.MkdirTemp("", "song_chunking")
	if err != nil {
		return &customerrors.ErrInternal{
			Message: "cannot create temp directory",
		}
	}
	defer os.RemoveAll(tempDir)

	songFileName := tempDir + "/" + song.Id + ".mp3"
	err = os.WriteFile(songFileName, input.File, 0777)
	if err != nil {
		return &customerrors.ErrInternal{
			Message: "cannot put song in the temp directory",
		}
	}

	cmd := exec.Command(
		"ffmpeg",
		"-i", songFileName,
		"-c:a", "libmp3lame",
		"-b:a", "128k",
		"-map", "0:0",
		"-f", "segment",
		"-segment_time", "5",
		"-segment_list", tempDir+"/"+song.Id+".m3u8",
		"-segment_format", "mpegts",
		tempDir+"/"+song.Id+"_"+"%d.ts",
	)
	err = cmd.Run()
	if err != nil {
		return &customerrors.ErrInvalidInput{
			Message: "cannot execute ffmpeg for song chunking",
		}
	}

	// Save header
	headerFile, err := os.ReadFile(tempDir + "/" + song.Id + ".m3u8")
	if err != nil {
		return &customerrors.ErrInternal{
			Message: "cannot read song header file",
		}
	}

	song.Length, err = utility.GetSongLengthFromM3U8(headerFile)
	if err != nil {
		return &customerrors.ErrInternal{
			Message: "cannot get song length from a m3u8",
		}
	}

	err = abstractions.SongChunkRepositoryInstance.
		AddSongChunk(song.Id, -1, headerFile)
	if err != nil {
		return err
	}

	// Save chunks
	for i := 0; true; i++ {
		headerFile, err := os.ReadFile(tempDir + "/" + song.Id + "_" + strconv.Itoa(i) + ".ts")
		if err != nil {
			break
		}
		err = abstractions.SongChunkRepositoryInstance.
			AddSongChunk(song.Id, i, headerFile)
		if err != nil {
			return err
		}
	}

	err = abstractions.SongRepositoryInstance.AddSong(song, albumCover, input.File)
	if err != nil {
		return err
	}

	return nil
}

func UpdateSongName(input models.UpdateSongNameInput) error {
	if check := utility.IsValidUUID(input.UserId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid user \"id\""}
	}
	if check := utility.IsValidUUID(input.SongId); !check {
		return &customerrors.ErrInvalidInput{
			Message: "invalid song \"id\""}
	}
	valid := utility.IsValidStructField(input, "Name")
	if !valid {
		return &customerrors.ErrInvalidInput{
			Message: "invalid \"new_name\" (must be 5-100 chars long)"}
	}

	song, err := abstractions.SongRepositoryInstance.
		GetSongInfo(input.SongId)
	if err != nil {
		return err
	}

	if song.CreatorId != input.UserId {
		return &customerrors.ErrInvalidInput{
			Message: "you are not the creator"}
	}

	err = abstractions.SongRepositoryInstance.
		UpdateSongName(input.SongId, input.Name)
	if err != nil {
		return err
	}
	return nil
}

func GetSongChunk(input models.GetSongChunkInput) ([]byte, error) {
	parts := strings.Split(input.FileName, ".")
	if len(parts) != 2 {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid file name"}
	}

	var resultFile []byte
	var err error
	if parts[1] == "m3u8" {
		if check := utility.IsValidUUID(parts[0]); !check {
			return nil, &customerrors.ErrInvalidInput{
				Message: "invalid song id"}
		}
		resultFile, err = abstractions.SongChunkRepositoryInstance.
			GetSongChunk(parts[0], -1)
		if err != nil {
			return nil, err
		}
		err = abstractions.SongRepositoryInstance.IncrementStreams(parts[0])
		if err != nil {
			return nil, err
		}
	} else if parts[1] == "ts" {
		leftparts := strings.Split(parts[0], "_")
		if check := utility.IsValidUUID(leftparts[0]); !check {
			return nil, &customerrors.ErrInvalidInput{
				Message: "invalid song id"}
		}
		index, err := strconv.Atoi(leftparts[1])
		if err != nil {
			return nil, &customerrors.ErrInvalidInput{
				Message: "invalid chunk id"}
		}
		resultFile, err = abstractions.SongChunkRepositoryInstance.
			GetSongChunk(leftparts[0], index)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, &customerrors.ErrInvalidInput{
			Message: "invalid file format"}
	}

	return resultFile, nil
}
