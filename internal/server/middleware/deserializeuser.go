package middleware

import (
	"spotigram/internal/server/abstractions"
	"spotigram/internal/server/config"
	"spotigram/internal/service/usecases"
	"spotigram/internal/utility"

	"github.com/gofiber/fiber/v2"
)

func DeserializeUser(ctx *fiber.Ctx) error {
	var access_token string

	if ctx.Cookies("access_token") != "" {
		access_token = ctx.Cookies("access_token")
	}

	if access_token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": "fail", "message": "you are not logged in"})
	}

	tokenClaims, err := utility.ValidateToken(
		access_token, config.Cfg.AccessToken.PublicKey)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	userUUID, err := abstractions.JWTCacheInstance.GetToken(tokenClaims.TokenUUID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "fail", "message": "token is invalid or session has expired"})
	}

	exists, err := usecases.DoesUserExist(userUUID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "message": "internal error"})
	} else if !exists {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

	ctx.Locals("user_uuid", userUUID)
	ctx.Locals("access_token_uuid", tokenClaims.TokenUUID)

	return ctx.Next()
}
