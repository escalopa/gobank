package api

import (
	"database/sql"
	"net/http"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/escalopa/gobank/util"
	"github.com/gin-gonic/gin"
)

type updateUserReq struct {
	FullName    string `json:"full_name" binding:"alpha"`
	Email       string `json:"email" binding:"email"`
	NewPassword string `json:"hashed_password"`
}

func (server *GinServer) updateUser(ctx *gin.Context) {
	var req updateUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, found := server.getUserIfExists(ctx, payload.Username)
	if !found {
		return
	}

	var params db.UpdateUserParams

	if req.Email == user.Email {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrEmailSameAsOld))
		return
	}

	if req.NewPassword != "" {
		hashedPassword, err := util.GenerateHashPassword(req.NewPassword)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		params.HashedPassword = sql.NullString{String: hashedPassword, Valid: true}
	}

	_, err := server.store.UpdateUser(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := fromUserToUserResponse(user)
	ctx.JSON(http.StatusOK, res)
}
