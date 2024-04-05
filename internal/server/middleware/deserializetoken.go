package middleware

import (
	"spotigram/internal/customerrors"
	"spotigram/internal/server/config"
	"spotigram/internal/service/models"
	"spotigram/internal/service/usecases"

	"github.com/gofiber/fiber/v2"
)

// A handler to deserialize access token.
func DeserializeTokenHandler(ctx *fiber.Ctx) error {
	input := models.DeserializeTokenInput{
		AccessToken: ctx.Cookies("access_token"),
	}
	userUuid, accessTokenUuid, err :=
		usecases.DeserializeToken(input, config.Cfg)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if errUnauthorized, ok := err.(*customerrors.ErrUnauthorized); ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "fail", "message": errUnauthorized.Error()})
		} else if errNotFound, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "fail", "message": errNotFound.Error()})
		}
	}

	ctx.Locals("user_uuid", userUuid)
	ctx.Locals("access_token_uuid", accessTokenUuid)

	return ctx.Next()
}
