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

// RenewAccessToken godoc
//
//	@Summary		renews an access token
//	@Description	renews an access token
//	@Tags			users
//	@Produce		json
//	@Param			body	body		renewAccessTokenReq	true	"Refresh token"
//	@Success		200		{object}	response.JSON{data=renewAccessTokenRes}
//	@Failure		400,500	{object}	response.JSON{}
//	@Router			/users/renew [post]
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
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, response.Err(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
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

	res := renewAccessTokenRes{AccessToken: accessToken, AccessTokenExpiresAt: accessPayload.ExpireAt}
	ctx.JSON(http.StatusAccepted, response.JSON{Data: res})
}
