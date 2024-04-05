package controllers

import (
	"encoding/json"

	"spotigram/internal/customerrors"

	"spotigram/internal/service/models"
	"spotigram/internal/service/usecases"

	"spotigram/internal/server/config"

	"github.com/gofiber/fiber/v2"
)

// A handler for user sign up.
func SignUpHandler(ctx *fiber.Ctx) error {
	sui := models.SignUpInput{}
	err := json.Unmarshal((ctx.Body()), &sui)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	err = usecases.SignUpUser(sui)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if errInvalidInput, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": errInvalidInput.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// A handler for user sign in.
func SignInHandler(ctx *fiber.Ctx) error {
	sii := models.SignInInput{}
	err := json.Unmarshal(ctx.Body(), &sii)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	userUuid, accessTokenDetails, refreshTokenDetails, err :=
		usecases.SignInUser(sii, config.Cfg)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if errInvalidInput, ok := err.(*customerrors.ErrInvalidInput); ok {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status": "fail", "message": errInvalidInput.Error()})
		} else if errNotFound, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": errNotFound.Error()})
		}
	}

	ctx.Cookie(&fiber.Cookie{
		Name:   "access_token",
		Value:  accessTokenDetails.Token,
		MaxAge: config.Cfg.AccessToken.MaxAge * 60,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:   "refresh_token",
		Value:  refreshTokenDetails.Token,
		MaxAge: config.Cfg.RefreshToken.MaxAge * 60,
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"uuid": userUuid})
}

// A handler for user logout.
func LogoutHandler(ctx *fiber.Ctx) error {
	input := models.LogoutInput{
		RefreshToken:    ctx.Cookies("refresh_token"),
		AccessTokenUUID: ctx.Locals("access_token_uuid").(string),
	}
	err := usecases.Logout(input, config.Cfg)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if errUnauthorized, ok := err.(*customerrors.ErrUnauthorized); ok {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status": "fail", "message": errUnauthorized.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

// A handler for access token refresh.
func RefreshAccessTokenHandler(ctx *fiber.Ctx) error {
	input := models.RefreshAccessTokenInput{
		RefreshToken: ctx.Cookies("refresh_token"),
	}
	accessTokenDetails, err :=
		usecases.RefreshAccessToken(input, config.Cfg)
	if err != nil {
		if errInternal, ok := err.(*customerrors.ErrInternal); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "fail", "message": errInternal.Error()})
		} else if errUnauthorized, ok := err.(*customerrors.ErrUnauthorized); ok {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status": "fail", "message": errUnauthorized.Error()})
		} else if errNotFound, ok := err.(*customerrors.ErrNotFound); ok {
			return ctx.Status(fiber.StatusNoContent).JSON(fiber.Map{
				"status": "fail", "message": errNotFound.Error()})
		}
	}

	ctx.Cookie(&fiber.Cookie{
		Name:   "access_token",
		Value:  accessTokenDetails.Token,
		MaxAge: config.Cfg.AccessToken.MaxAge * 60,
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessTokenDetails.Token})
}
