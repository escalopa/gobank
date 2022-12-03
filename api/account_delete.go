package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type deleteAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *GinServer) deleteAccounts(ctx *gin.Context) {
	var req deleteAccountReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, isValid := server.isValidAccount(ctx, req.ID)
	if !isValid {
		return
	}

	if !isUserAccountOwner(ctx, account) {
		ctx.JSON(http.StatusUnauthorized, ErrNotAccountOwner)
		return
	}

	err := server.store.DeleteAccount(ctx, req.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"id": req.ID})
}
