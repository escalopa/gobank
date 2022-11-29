package api

import (
	"fmt"
	"net/http"

	db "github.com/escalopa/go-bank/db/sqlc"
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

func (server *Server) createTransfer(ctx *gin.Context) {
	var req createTransferReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	if !server.validateTransferCurrency(ctx, req.FromAccountID, req.ToAccountID) {
		return
	}

	if req.FromAccountID == req.ToAccountID {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf(fmt.Sprintf("can't transfer to the same account, req.FromAccountId=%d, req.ToAccount=%d",
			req.FromAccountID,
			req.ToAccountID,
		))))
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

func (server *Server) validateTransferCurrency(ctx *gin.Context, accountID1, accountID2 int64) bool {

	account1, err := server.store.GetAccount(ctx, accountID1)
	if err != nil {
		server.handleGetDataBaseError(ctx, err)
		return false
	}

	account2, err := server.store.GetAccount(ctx, accountID2)
	if err != nil {
		server.handleGetDataBaseError(ctx, err)
		return false
	}

	if account1.Currency != account2.Currency {
		ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf(fmt.Sprintf("currency mismatch account1.currency=%s, account2.currency=%s",
			account1.Currency,
			account2.Currency))))
		return false
	}

	return true
}
