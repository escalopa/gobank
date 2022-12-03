package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type getAccountReq struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *GinServer) getAccount(ctx *gin.Context) {
	var req getAccountReq
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

	ctx.JSON(http.StatusOK, account)
}
