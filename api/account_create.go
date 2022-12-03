package api

import (
	"database/sql"
	"net/http"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountReq struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *GinServer) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	account, err := server.store.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    payload.Username,
		Balance:  0,
		Currency: req.Currency,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

func isUserAccountOwner(ctx *gin.Context, account db.Account) bool {
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	return payload.Username == account.Owner
}

func (server *GinServer) isValidAccount(ctx *gin.Context, accountID int64) (db.Account, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(err))
				return db.Account{}, false
			}

			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return db.Account{}, false
		}
	}
	return account, true
}
