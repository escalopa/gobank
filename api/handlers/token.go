package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/escalopa/gobank/api/handlers/response"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type renewAccessTokenRes struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_expires_at"`
}

func (s *GinServer) renewAccessToken(ctx *gin.Context) {
	var req renewAccessTokenReq
	if err := parseBody(ctx, &req); err != nil {
		return
	}

	refreshPayload, err := s.tm.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	session, err := s.db.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, response.Err(err))
				return
			}

			ctx.JSON(http.StatusInternalServerError, response.Err(err))
			return
		}
	}

	if session.IsBlocked {
		ctx.JSON(http.StatusUnauthorized, response.Err(ErrBlockedRefreshToken))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, response.Err(ErrExpiredRefreshToken))
		return
	}

	if session.Username != refreshPayload.Username {
		ctx.JSON(http.StatusUnauthorized, response.Err(ErrMismatchedRefreshTokens))
		return
	}

	if session.RefreshToken != req.RefreshToken {
		ctx.JSON(http.StatusUnauthorized, response.Err(ErrMismatchedRefreshTokens))
		return
	}

	// Generate New Access Token for User
	accessToken, accessPayload, err := s.tm.CreateToken(session.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	res := renewAccessTokenRes{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpireAt,
	}
	ctx.JSON(http.StatusAccepted, res)
}
