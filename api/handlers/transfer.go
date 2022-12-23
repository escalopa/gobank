package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/escalopa/gobank/api/handlers/response"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferResponse struct {
	ID          int64      `json:"id"`
	FromAccount db.Account `json:"from_account"`
	FromEntry   db.Entry   `json:"from_entry"`
	ToAccountID int64      `json:"to_account_id"`
	Amount      int64      `json:"amount"`
	CreatedAt   time.Time  `json:"created_at"`
}

type createTransferReq struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gte=1"`
}

// CreateTransfer godoc
//
//	@Summary		creates a new transfer between two accounts
//	@Description	creates a new transfer between two accounts
//	@Tags			transfers
//	@Accept			json
//	@Produce		json
//	@Param			body	body		createTransferReq	true	"Transfer to create"
//	@Success		200		{object}	response.JSON{data=transferResponse}
//	@Failure		400,500	{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/transfers [post]
func (s *GinServer) createTransfer(ctx *gin.Context) {
	var req createTransferReq
	if err := parseBody(ctx, &req); err != nil {
		return
	}

	fromAccount, to, isValid := s.validateTransfer(ctx, req.FromAccountID, req.ToAccountID)
	if !isValid {
		return
	}

	_ = to
	if !isUserAccountOwner(ctx, fromAccount) {
		ctx.JSON(http.StatusUnauthorized, ErrNotAccountOwner)
		return
	}

	arg := db.TransferTxParam{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := s.db.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	res := fromTransferTxToTransferResponse(result)
	ctx.JSON(http.StatusOK, res)

}

func (s *GinServer) validateTransfer(ctx *gin.Context, accountID1, accountID2 int64) (from db.Account, to db.Account, isValid bool) {
	if accountID1 == accountID2 {
		ctx.JSON(http.StatusBadRequest, response.Err(ErrSameAccountTransfer(accountID1, accountID2)))
		return
	}

	from, isValid = s.isValidAccount(ctx, accountID1)
	if !isValid {
		return
	}

	to, isValid = s.isValidAccount(ctx, accountID2)
	if !isValid {
		return
	}

	isValid = from.Currency == to.Currency
	if !isValid {
		ctx.JSON(http.StatusBadRequest, response.Err(ErrCurrencyMismatch(from.Currency, to.Currency)))
		return
	}

	if from.IsDeleted {
		isValid = false
		ctx.JSON(http.StatusBadRequest, response.Err(ErrAccountDeleted(from.ID)))
		return
	}

	if to.IsDeleted {
		isValid = false
		ctx.JSON(http.StatusBadRequest, response.Err(ErrAccountDeleted(to.ID)))
		return
	}

	return
}

type getTransferReq struct {
	AccountID int64 `uri:"id" binding:"required,min=1"`
}

// GetTransfers godoc
//
//	@Summary		gets all transfers for an account
//	@Description	gets all transfers for an account
//	@Tags			transfers
//	@Accept			json
//	@Produce		json
//	@Param			id			path		int64	true	"Account ID"
//	@Param			page_id		query		int32	true	"Page ID"
//	@Param			page_size	query		int32	true	"Page Size"
//	@Success		200			{object}	response.JSON{data=transferResponse}
//	@Failure		400,500		{object}	response.JSON{}
//	@Security		bearerAuth
//	@Router			/transfers/{id} [get]
func (s *GinServer) getTransfers(ctx *gin.Context) {
	var req getTransferReq
	var pgQuery *paginationQuery
	var err error

	if err := parseUri(ctx, &req); err != nil {
		return
	}

	if pgQuery, err = parsePagination(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err(err))
		return
	}

	account, ok := s.isValidAccount(ctx, req.AccountID)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.Err(errors.New("invalid account id")))
		return
	}

	if !isUserAccountOwner(ctx, account) {
		ctx.JSON(http.StatusUnauthorized, response.Err(ErrNotAccountOwner))
		return
	}

	transfers, err := s.db.ListTransfers(ctx, db.ListTransfersParams{
		AccountID: req.AccountID, PageSize: pgQuery.Limit, PageID: (pgQuery.Offset - 1) * pgQuery.Offset,
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Err(err))
		return
	}

	var responsesTransfer []*transferResponse
	for _, transfer := range transfers {
		responsesTransfer = append(responsesTransfer, mapTransferToResponse(transfer))
	}

	ctx.JSON(http.StatusOK, response.Success(responsesTransfer))
}
