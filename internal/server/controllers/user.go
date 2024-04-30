package controllers

import (
	"encoding/json"

	"spotigram/internal/customerrors"
	"spotigram/internal/service/models"
	"spotigram/internal/service/usecases"

	"github.com/gofiber/fiber/v2"
)

// A handler to send user list.
func UsersHandler(ctx *fiber.Ctx) error {
	input := models.GetUsersInput{}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	users, err := usecases.GetUsers(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if _, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": "users not found"})
		} else if _, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": "invalid offset"})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(users)
}

// A handler to send a user's info.
func UserInfoHandler(ctx *fiber.Ctx) error {
	uii := models.GetUserInfoInput{}
	err := json.Unmarshal((ctx.Body()), &uii)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	user, err := usecases.GetUserInfo(uii)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if _, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": "user not found"})
		} else if _, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": "invalid uuid"})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":       user.Id,
		"email":    user.Email,
		"name":     user.Name,
		"verified": user.Verified,
	})
}

// A handler to send current user public key.
func UserPublicKeyHandler(ctx *fiber.Ctx) error {

	input := models.GetPublicKeyInput{}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	key, err := usecases.GetPublicKey(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if _, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": "user not found"})
		} else if _, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": "invalid uuid"})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"public_key": key,
	})
}

// A handler to send current user's picture.
func UserPictureHandler(ctx *fiber.Ctx) error {
	input := models.GetPictureInput{}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	pic, err := usecases.GetPicture(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if _, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": "user not found"})
		} else if _, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": "invalid uuid"})
		}
	}

	return ctx.Status(fiber.StatusOK).Send(pic)
}
