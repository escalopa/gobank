package handlers

import (
	"net/http"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type transferResponse struct {
	ID          int64      `json:"id"`
	FromAccount db.Account `json:"from_account"`
	FromEntry   db.Entry   `json:"from_entry"`
	ToAccountID int64      `json:"to_account_id"`
	Amount      int64      `json:"amount"`
}

type createTransferReq struct {
	FromAccountID int64 `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64 `json:"to_account_id" binding:"required,min=1"`
	Amount        int64 `json:"amount" binding:"required,gte=1"`
}

func (server *GinServer) createTransfer(ctx *gin.Context) {
	var req createTransferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	fromAccount, to, isValid := server.validateTransfer(ctx, req.FromAccountID, req.ToAccountID)
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

	result, err := server.store.TransferTx(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := fromTransferTxToTransferResponse(result)
	ctx.JSON(http.StatusOK, res)

}

func (server *GinServer) validateTransfer(ctx *gin.Context, accountID1, accountID2 int64) (from db.Account, to db.Account, isValid bool) {
	if accountID1 == accountID2 {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrSameAccountTransfer(accountID1, accountID2)))
		return
	}

	from, isValid = server.isValidAccount(ctx, accountID1)
	if !isValid {
		return
	}

	to, isValid = server.isValidAccount(ctx, accountID2)
	if !isValid {
		return
	}

	isValid = from.Currency == to.Currency
	if !isValid {
		ctx.JSON(http.StatusBadRequest, errorResponse(ErrCurrencyMismatch(from.Currency, to.Currency)))
		return
	}

	return
}
