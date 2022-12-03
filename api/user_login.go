package api

import (
	"net/http"
	"time"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type loginUserReq struct {
	Username string `json:"username" binding:"required,min=6,max=16,alphanum"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

type loginUserRes struct {
	SessionID             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_expires_at"`
	User                  userResponse `json:"user"`
}

func (server *GinServer) loginUser(ctx *gin.Context) {
	var req loginUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// Get User from DB by Username
	user, found := server.getUserIfExists(ctx, req.Username)
	if !found {
		return
	}

	// Check User's Password
	err := util.CheckHashedPassword(user.HashedPassword, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	// Generate New Access Token for User
	accessToken, accessPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenExpiration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Generate New Access Token for User
	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(user.Username, server.config.RefreshTokenExpiration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// Create new session for User
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     req.Username,
		RefreshToken: refreshToken,
		UserAgent:    ctx.Request.UserAgent(),
		ClientIp:     ctx.ClientIP(),
		ExpiresAt:    refreshPayload.ExpireAt,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := loginUserRes{
		SessionID:             session.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  accessPayload.ExpireAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpireAt,
		User:                  fromUserToUserResponse(user),
	}
	ctx.JSON(http.StatusAccepted, res)
}
