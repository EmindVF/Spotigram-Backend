package controllers

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/service/usecases"

	"github.com/gofiber/fiber/v2"
)

func MyInfoHandler(ctx *fiber.Ctx) error {
	user, err := usecases.GetUser(ctx.Locals("user_uuid").(string))
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
