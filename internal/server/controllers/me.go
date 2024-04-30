package controllers

import (
	"encoding/json"

	"spotigram/internal/customerrors"
	"spotigram/internal/service/models"
	"spotigram/internal/service/usecases"

	"github.com/gofiber/fiber/v2"
)

// A handler to send current user info.
func MyInfoHandler(ctx *fiber.Ctx) error {
	input := models.GetUserInfoInput{
		UserUUID: ctx.Locals("user_uuid").(string),
	}
	user, err := usecases.GetUserInfo(input)
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

// A handler to send current user friends.
func FriendsHandler(ctx *fiber.Ctx) error {
	input := models.GetFriendsInput{
		UserUUID: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	friends, err := usecases.GetFriends(input)
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

	return ctx.Status(fiber.StatusOK).JSON(friends)
}

// A handler to send current user friend requests sent.
func FriendRequestSentHandler(ctx *fiber.Ctx) error {
	input := models.GetFriendRequestsSentInput{
		UserUUID: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	friendRequests, err := usecases.GetFriendRequestsSent(input)
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

	return ctx.Status(fiber.StatusOK).JSON(friendRequests)
}

// A handler to send current user friend requests received.
func FriendRequestReceivedHandler(ctx *fiber.Ctx) error {
	input := models.GetFriendRequestsReceivedInput{
		UserUUID: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	friendRequests, err := usecases.GetFriendRequestsReceived(input)
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

	return ctx.Status(fiber.StatusOK).JSON(friendRequests)
}

// A handler to send current user info.
func MyPublicKeyHandler(ctx *fiber.Ctx) error {
	input := models.GetPublicKeyInput{
		UserUUID: ctx.Locals("user_uuid").(string),
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
func MyPictureHandler(ctx *fiber.Ctx) error {
	input := models.GetPictureInput{
		UserUUID: ctx.Locals("user_uuid").(string),
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

// A handler to change user's name.
func ChangeNameHandler(ctx *fiber.Ctx) error {

	input := models.ChangeNameInput{}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	input.UserUUID = ctx.Locals("user_uuid").(string)
	err = usecases.ChangeName(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if _, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": "user not found"})
		} else if errInvalidInput, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": errInvalidInput.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// A handler to change user's password.
func ChangePasswordHandler(ctx *fiber.Ctx) error {

	input := models.ChangePasswordInput{}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	input.UserUUID = ctx.Locals("user_uuid").(string)
	err = usecases.ChangePassword(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if _, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": "user not found"})
		} else if errInvalidInput, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": errInvalidInput.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// A handler to change user's public key for end-to-end encryption.
func ChangePublicKeyHandler(ctx *fiber.Ctx) error {

	input := models.ChangePublicKeyInput{}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	input.UserUUID = ctx.Locals("user_uuid").(string)
	err = usecases.ChangePublicKey(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if _, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": "user not found"})
		} else if errInvalidInput, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": errInvalidInput.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// A handler to change user's picture.
func ChangePictureHandler(ctx *fiber.Ctx) error {

	input := models.ChangePictureInput{
		Image: ctx.Body(),
	}

	if input.Image == nil || len(input.Image) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid image, must be a png or a jpg (under 5 megabytes)"})
	}

	input.UserUUID = ctx.Locals("user_uuid").(string)
	err := usecases.ChangePicture(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if _, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": "user not found"})
		} else if errInvalidInput, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": errInvalidInput.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}
