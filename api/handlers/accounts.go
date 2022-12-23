package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/escalopa/gobank/api/handlers/response"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type accountResponse struct {
	ID        int64     `json:"id"`
	Balance   int64     `json:"balance"`
	Currency  string    `json:"currency"`
	CreatedAt time.Time `json:"created_at"`
}

type createAccountReq struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// CreateAccount godoc
//
//	@Summary		creates a new account for the currently logged-in user
//	@Description	creates a new account for the currently logged-in user
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createAccountReq	true	"Account to create"
//	@Success		200		{object}	response.JSON{data=accountResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts [post]
func (s *GinServer) createAccount(ctx *gin.Context) {
	var req createAccountReq
	if err := parseBody(ctx, &req); err != nil {
		return
	}

	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	account, err := s.db.CreateAccount(ctx, db.CreateAccountParams{
		Owner:    payload.Username,
		Balance:  1000,
		Currency: req.Currency,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, response.Err(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	ctx.JSON(http.StatusCreated, response.Success(mapAccountToResponse(account)))
}

func isUserAccountOwner(ctx *gin.Context, account db.Account) bool {
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	log.Println(payload.Username, account.Owner)
	return payload.Username == account.Owner
}

func (s *GinServer) isValidAccount(ctx *gin.Context, accountID int64) (db.Account, bool) {
	account, err := s.db.GetAccount(ctx, accountID)
	if err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				ctx.JSON(http.StatusNotFound, response.Err(err))
				return db.Account{}, false
			}

			ctx.JSON(http.StatusInternalServerError, response.Err(err))
			return db.Account{}, false
		}
	}
	return account, true
}

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// GetAccount godoc
//
//	@Summary		gets an account by id
//	@Description	gets an account by id
//	@Tags			accounts
//	@Produce		json
//	@Param			id		path		int64	true	"Account ID"
//	@Success		200		{object}	response.JSON{data=accountResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts/{id} [get]
func (s *GinServer) getAccount(ctx *gin.Context) {
	var req getAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err(err))
		return
	}

	account, isValid := s.isValidAccount(ctx, req.ID)
	if !isValid {
		return
	}

	if !isUserAccountOwner(ctx, account) {
		ctx.JSON(http.StatusUnauthorized, ErrNotAccountOwner)
		return
	}

	ctx.JSON(http.StatusOK, response.Success(mapAccountToResponse(account)))
}

// GetDeletedAccounts godoc
//
//	@Summary
//	@Description	gets a list of accounts for the currently logged-in user
//	@Tags			accounts
//	@Produce		json
//	@Success		200		{object}	response.JSON{data=[]accountResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts/del [get]
func (s *GinServer) getDeletedAccounts(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	accounts, err := s.db.GetDeletedAccounts(ctx, payload.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	var resp []*accountResponse
	for _, account := range accounts {
		resp = append(resp, mapAccountToResponse(account))
	}

	ctx.JSON(http.StatusOK, response.Success(resp))
}

// GetAccounts godoc
//
//	@Summary		gets a list of accounts for the currently logged-in user
//	@Description	gets a list of accounts for the currently logged-in user
//	@Tags			accounts
//	@Produce		json
//	@Success		200		{object}	response.JSON{data=[]accountResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts [get]
func (s *GinServer) getAccounts(ctx *gin.Context) {
	payload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	accounts, err := s.db.GetAccounts(ctx, payload.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	var accountResponses []*accountResponse
	for _, account := range accounts {
		accountResponses = append(accountResponses, mapAccountToResponse(account))
	}

	ctx.JSON(http.StatusOK, response.Success(accountResponses))
}

type deleteAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// DeleteAccount godoc
//
//	@Summary		deletes an account by id for the currently logged-in user
//	@Description	deletes an account by id for the currently logged-in user
//	@Tags			accounts
//	@Produce		json
//	@Param			id		path		int64	true	"Account ID"
//	@Success		200		{object}	response.JSON{data=int64}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts/{id} [delete]
func (s *GinServer) deleteAccount(ctx *gin.Context) {
	var req deleteAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err(err))
		return
	}

	account, isValid := s.isValidAccount(ctx, req.ID)
	if !isValid {
		return
	}

	if !isUserAccountOwner(ctx, account) {
		ctx.JSON(http.StatusUnauthorized, response.Err(ErrNotAccountOwner))
		return
	}

	err := s.db.DeleteAccount(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(account.ID))
}

type restoreAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// RestoreAccount godoc
//
//	@Summary		deletes an account by id for the currently logged-in user
//	@Description	deletes an account by id for the currently logged-in user
//	@Tags			accounts
//	@Produce		json
//	@Param			id		path		int64	true	"Account ID"
//	@Success		200		{object}	response.JSON{data=int64}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/accounts/res/{id} [patch]
func (s *GinServer) restoreAccount(ctx *gin.Context) {
	var req restoreAccountReq
	if err := parseUri(ctx, &req); err != nil {
		return
	}

	account, isValid := s.isValidAccount(ctx, req.ID)
	if !isValid {
		return
	}

	if !isUserAccountOwner(ctx, account) {
		ctx.JSON(http.StatusUnauthorized, response.Err(ErrNotAccountOwner))
		return
	}

	err := s.db.RestoreAccount(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Success(account.ID))
}
