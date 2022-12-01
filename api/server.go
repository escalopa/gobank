package api

import (
	"fmt"

	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/escalopa/gobank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create tokenMaker, %w", err)
	}

	server := &Server{config: config, tokenMaker: maker, store: store}

	server.setupValidator()
	server.setupRouter()
	return server, nil
}

func (server *Server) Start(address string) error {
	if err := server.router.Run(address); err != nil {
		return err
	}
	return nil
}

func (server *Server) setupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
}

func (server *Server) setupRouter() {
	router := gin.Default()

	authGroup := router.Group("/").Use(authMiddleware(server.tokenMaker))

	{
		// Account Routing
		authGroup.POST("/api/accounts", server.createAccount)
		authGroup.GET("/api/accounts/:id", server.getAccount)
		authGroup.GET("/api/accounts", server.listAccounts)
		authGroup.DELETE("/api/accounts/:id", server.deleteAccounts)

		// Transfer Routing
		authGroup.POST("/api/transfers", server.createTransfer)

		// User Routing
		authGroup.GET("api/users/:username", server.getUser)
	}

	// Authenticate Routing
	router.POST("api/users", server.createUser)
	router.POST("api/users/login", server.loginUser)

	server.router = router
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
