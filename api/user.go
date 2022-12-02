package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/escalopa/gobank/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

type createUserReq struct {
	Username        string `json:"username"  binding:"required,min=6,max=16,alphanum"`
	FullName        string `json:"full_name" binding:"required"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6,max=16"`
	PasswordConfirm string `json:"password_confirm" binding:"required,eqfield=Password"`
}

func (server *GinServer) createUser(ctx *gin.Context) {
	var req createUserReq
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashPassword, err := util.GenerateHashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("couldn't hash password, err: %v", err))
		return
	}

	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := fromUserToUserResponse(user)
	ctx.JSON(http.StatusCreated, res)
}

func (server *GinServer) getUser(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, isValid := server.isValidUser(ctx, payload.Username)
	if !isValid {
		return
	}

	res := fromUserToUserResponse(user)
	ctx.JSON(http.StatusOK, res)
}

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
	user, isValid := server.isValidUser(ctx, req.Username)
	if !isValid {
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

func (server *GinServer) isValidUser(ctx *gin.Context, username string) (db.User, bool) {
	user, err := server.store.GetUser(ctx, username)
	if err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return db.User{}, false
			}

			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return db.User{}, false
		}
	}
	return user, true
}

// func (server *GinServer) getCurrentUser(ctx *gin.Context) db.User {
// 	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
// 	user, isValid := server.isValidUser(ctx, payload.Username)
// 	if !isValid {
// 		return db.User{}
// 	}
// 	return user
// }
