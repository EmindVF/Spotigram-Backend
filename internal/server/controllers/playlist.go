package controllers

import (
	"encoding/json"
	"spotigram/internal/customerrors"
	"spotigram/internal/service/models"
	"spotigram/internal/service/usecases"

	"github.com/gofiber/fiber/v2"
)

// A handler to send current user info.
func PlaylistsHandler(ctx *fiber.Ctx) error {
	input := models.GetPlaylistsInput{
		UserId: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	playlists, err := usecases.GetPlaylists(input)
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

	return ctx.Status(fiber.StatusOK).JSON(playlists)
}

// A handler to send current user info.
func PlaylistSongsHandler(ctx *fiber.Ctx) error {
	input := models.GetPlaylistSongsInput{
		UserId: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	songs, err := usecases.GetPlaylistSongs(input)
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

// A handler to send current user info.
func DeletePlaylistHandler(ctx *fiber.Ctx) error {
	input := models.DeletePlaylistInput{
		UserId: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	err = usecases.DeletePlaylist(input)
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

// A handler to send current user info.
func AddPlaylistHandler(ctx *fiber.Ctx) error {
	input := models.AddPlaylistInput{
		UserId: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	uuid, err := usecases.AddPlaylist(input)
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

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"id": uuid})
}

// A handler to send current user info.
func AddPlaylistSongHandler(ctx *fiber.Ctx) error {
	input := models.AddPlaylistSongInput{
		UserId: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	err = usecases.AddPlaylistSong(input)
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

// A handler to send current user info.
func DeletePlaylistSongHandler(ctx *fiber.Ctx) error {
	input := models.DeletePlaylistSongInput{
		UserId: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	err = usecases.DeletePlaylistSong(input)
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
