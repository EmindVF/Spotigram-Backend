package controllers

import (
	"encoding/json"
	"spotigram/internal/customerrors"
	"spotigram/internal/service/models"
	"spotigram/internal/service/usecases"

	"github.com/gofiber/fiber/v2"
)

// A handler to send current user info.
func SongsHandler(ctx *fiber.Ctx) error {
	input := models.GetSongsInput{}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	songs, err := usecases.GetSongs(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if newErr, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		} else if newErr, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(songs)
}

func SongInfoHandler(ctx *fiber.Ctx) error {
	input := models.GetSongInfoInput{}
	err := json.Unmarshal(ctx.Body(), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	song, err := usecases.GetSongInfo(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if newErr, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		} else if newErr, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(song)
}

func SongPictureHandler(ctx *fiber.Ctx) error {
	input := models.GetSongPictureInput{}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	pic, err := usecases.GetSongPicture(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if _, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": "song not found"})
		} else if _, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": "invalid uuid"})
		}
	}

	return ctx.Status(fiber.StatusOK).Send(pic)
}

func RenameSongHandler(ctx *fiber.Ctx) error {
	input := models.UpdateSongNameInput{
		UserId: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	err = usecases.UpdateSongName(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if _, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": "song not found"})
		} else if _, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": "invalid uuid"})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func UploadSongHandler(ctx *fiber.Ctx) error {
	input := models.AddSongInput{
		UserId: ctx.Locals("user_uuid").(string),
		Name:   ctx.Params("songname"),
		File:   ctx.Body(),
	}
	err := usecases.AddSong(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if newErr, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		} else if newErr, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func DeleteSongHandler(ctx *fiber.Ctx) error {
	input := models.DeleteSongInput{
		UserId: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	err = usecases.DeleteSong(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if newErr, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		} else if newErr, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		} else if newErr, ok := err.(*customerrors.ErrUnauthorized); ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func DownloadSongHandler(ctx *fiber.Ctx) error {
	input := models.GetSongFileInput{}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	file, err := usecases.GetSongFile(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if newErr, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		} else if newErr, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).Send(file)
}

func GetSongChunk(ctx *fiber.Ctx) error {
	input := models.GetSongChunkInput{
		FileName: ctx.Params("filename"),
	}

	file, err := usecases.GetSongChunk(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if newErr, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		} else if newErr, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": newErr.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).Send(file)
}
