package handlers

import (
	"net/http"

	"github.com/escalopa/gobank/api/handlers/response"

	"github.com/gin-gonic/gin"
)

func parseBody(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindJSON(obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err(err))
		return err
	}
	return nil
}

func parseUri(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindUri(obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err(err))
		return err
	}
	return nil
}

func parseQuery(ctx *gin.Context, obj interface{}) error {
	err := ctx.ShouldBindQuery(obj)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.Err(err))
		return err
	}
	return nil
}
