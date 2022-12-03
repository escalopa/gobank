package api

import (
	"database/sql"
	"net/http"
	"time"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/gin-gonic/gin"
)

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	CreatedAt         time.Time `json:"created_at"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

func (server *GinServer) getUser(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	user, found := server.getUserIfExists(ctx, payload.Username)
	if !found {
		return
	}

	res := fromUserToUserResponse(user)
	ctx.JSON(http.StatusOK, res)
}

func (server *GinServer) getUserIfExists(ctx *gin.Context, username string) (db.User, bool) {
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
