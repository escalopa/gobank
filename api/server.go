package api

import (
	db "github.com/escalopa/go-bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{store: store}

	router := gin.Default()

	// Account Routing
	router.POST("/api/accounts", server.createAccount)
	router.GET("/api/accounts/:id", server.getAccount)
	router.GET("/api/accounts", server.listAccounts)
	router.DELETE("/api/accounts/:id", server.deleteAccounts)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	if err := server.router.Run(address); err != nil {
		return err
	}
	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
