package controllers

import (
	"encoding/json"
	"spotigram/internal/customerrors"
	"spotigram/internal/service/models"
	"spotigram/internal/service/usecases"

	"github.com/gofiber/fiber/v2"
)

// A handler to send all chat messages.
func ChatMessagesHandler(ctx *fiber.Ctx) error {
	input := models.GetMessagesInput{
		UserId: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	messages, err := usecases.GetMessages(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if errNotFound, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": errNotFound.Error()})
		} else if errInvalidInput, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": errInvalidInput.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(messages)
}

// A handler to send unread messages of a chat
func ChatUnreadMessagesHandler(ctx *fiber.Ctx) error {
	input := models.GetUnreadMessagesInput{
		UserId: ctx.Locals("user_uuid").(string),
	}
	err := json.Unmarshal((ctx.Body()), &input)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}
	messages, err := usecases.GetUnreadMessages(input)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if errNotFound, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": errNotFound.Error()})
		} else if errInvalidInput, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": errInvalidInput.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(messages)
}
