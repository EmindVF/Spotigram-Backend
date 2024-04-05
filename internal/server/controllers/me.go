package controllers

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/service/models"
	"spotigram/internal/service/usecases"

	"github.com/gofiber/fiber/v2"
)

// A handler to send current user info.
func MyInfoHandler(ctx *fiber.Ctx) error {
	input := models.MyInfoInput{
		AccesssTokenUUID: ctx.Locals("user_uuid").(string),
	}
	user, err := usecases.MyInfo(input)
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
