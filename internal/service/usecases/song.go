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
	songs, err :=
		abstractions.SongRepositoryInstance.GetSongs(gsi.Offset)
	if err != nil {
		return nil, err
	}
	return songs, nil
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

	err := abstractions.SongRepositoryInstance.
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

	if !utility.IsValidMP3(input.File) {
		return &customerrors.ErrInvalidInput{
			Message: "file is not an mp3"}
	}

	songLen, err := utility.GetMP3Length(input.File)
	if err != nil {
		return &customerrors.ErrInternal{
			Message: "cannot determine song length",
		}
	}

	pic, _, err := utility.GetMP3AlbumCover(input.File)
	if err != nil {
		return &customerrors.ErrInternal{
			Message: "cannot read tags of the song",
		}
	}

	var albumCover []byte

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
	} else {
		albumCover = nil
	}

	song := models.Song{
		Id:        utility.GenerateUUID(),
		CreatorId: input.UserId,
		Name:      input.Name,
		Length:    songLen,
	}

	err = abstractions.SongRepositoryInstance.AddSong(song, albumCover, input.File)
	if err != nil {
		return err
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
		return &customerrors.ErrInternal{
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
