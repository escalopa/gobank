package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/escalopa/go-bank/db/sqlc"
	"github.com/escalopa/go-bank/token"
	"github.com/escalopa/go-bank/util"
	"github.com/gin-gonic/gin"
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

func (server *Server) createUser(ctx *gin.Context) {
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

func (server *Server) getUser(ctx *gin.Context) {
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
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
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
	accessToken, err := server.tokenMaker.CreateToken(user.Username, server.config.App.TokenExpiration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := loginUserRes{AccessToken: accessToken, User: fromUserToUserResponse(user)}
	ctx.JSON(http.StatusAccepted, res)
}

func (server *Server) isValidUser(ctx *gin.Context, username string) (db.User, bool) {
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

// func (server *Server) getCurrentUser(ctx *gin.Context) db.User {
// 	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
// 	user, isValid := server.isValidUser(ctx, payload.Username)
// 	if !isValid {
// 		return db.User{}
// 	}
// 	return user
// }
