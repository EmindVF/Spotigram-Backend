package controllers

import (
	"encoding/json"
	"time"

	"spotigram/internal/customerrors"
	"spotigram/internal/service/models"
	"spotigram/internal/service/usecases"

	"spotigram/internal/server/abstractions"
	"spotigram/internal/server/config"

	"spotigram/internal/utility"

	"github.com/gofiber/fiber/v2"
)

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

func SignInHandler(ctx *fiber.Ctx) error {
	sii := models.SignInInput{}
	err := json.Unmarshal(ctx.Body(), &sii)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	userUuid, err := usecases.SignInUser(sii)
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

	accessTokenDetails, err := utility.CreateToken(
		userUuid, config.Cfg.AccessToken.ExpiresIn, config.Cfg.AccessToken.PrivateKey)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	refreshTokenDetails, err := utility.CreateToken(
		userUuid, config.Cfg.RefreshToken.ExpiresIn, config.Cfg.RefreshToken.PrivateKey)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	now := time.Now()

	err = abstractions.JWTCacheInstance.SetToken(
		accessTokenDetails.TokenUUID,
		userUuid,
		(time.Unix(accessTokenDetails.ExpiresIn, 0).Sub(now)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	err = abstractions.JWTCacheInstance.SetToken(
		refreshTokenDetails.TokenUUID,
		userUuid,
		(time.Unix(refreshTokenDetails.ExpiresIn, 0).Sub(now)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
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

func LogoutHandler(ctx *fiber.Ctx) error {
	refresh_token := ctx.Cookies("refresh_token")
	if refresh_token == "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "fail", "message": "token is invalid or session has expired"})
	}

	refreshTokenClaims, err := utility.ValidateToken(refresh_token, config.Cfg.RefreshToken.PublicKey)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	accessTokenUUID := ctx.Locals("access_token_uuid").(string)
	_, err = abstractions.JWTCacheInstance.DeleteRefreshAndAccessToken(
		refreshTokenClaims.TokenUUID, accessTokenUUID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
}

func RefreshAccessTokenHandler(ctx *fiber.Ctx) error {
	refresh_token := ctx.Cookies("refresh_token")

	if refresh_token == "" {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "fail", "message": "could not refresh access token"})
	}

	tokenClaims, err := utility.ValidateToken(refresh_token, config.Cfg.RefreshToken.PublicKey)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	userUuid, err := abstractions.JWTCacheInstance.GetToken(tokenClaims.TokenUUID)
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "fail", "message": "could not refresh access token"})
	}

	exists, err := usecases.DoesUserExist(userUuid)
	if !exists {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status": "fail", "message": "the user belonging to this token no logger exists"})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	accessTokenDetails, err := utility.CreateToken(
		userUuid, config.Cfg.AccessToken.ExpiresIn, config.Cfg.AccessToken.PrivateKey)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	now := time.Now()

	err = abstractions.JWTCacheInstance.SetToken(
		accessTokenDetails.TokenUUID,
		userUuid,
		(time.Unix(accessTokenDetails.ExpiresIn, 0).Sub(now)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": "fail", "message": err.Error()})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:   "access_token",
		Value:  accessTokenDetails.Token,
		MaxAge: config.Cfg.AccessToken.MaxAge * 60,
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"access_token": accessTokenDetails.Token})
}
