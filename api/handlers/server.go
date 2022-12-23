package handlers

import (
	"fmt"

	_ "github.com/escalopa/gobank/api/docs"
	db "github.com/escalopa/gobank/db/sqlc"
	"github.com/escalopa/gobank/token"
	"github.com/escalopa/gobank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type GinServer struct {
	config *util.Config
	db     db.Store
	tm     token.Maker
	router *gin.Engine
}

func NewServer(config *util.Config, store db.Store) (*GinServer, error) {
	maker, err := token.NewPasetoMaker(config.Get("SYMMETRIC_KEY"))
	if err != nil {
		return nil, fmt.Errorf("cannot create tokenMaker, %w", err)
	}

	s := &GinServer{config: config, tm: maker, db: store}

	gin.SetMode(gin.ReleaseMode)
	s.setupValidator()
	s.setupRouter()
	s.setupSwagger()
	return s, nil
}

func (s *GinServer) Start(address string) error {
	if err := s.router.Run(address); err != nil {
		return err
	}
	return nil

}

func (s *GinServer) setupValidator() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
}

func (s *GinServer) setupRouter() {
	router := gin.Default()

	auth := router.Group("/").Use(authMiddleware(s.tm))

	// Account Routes
	auth.POST("/api/accounts", s.createAccount)
	auth.GET("/api/accounts/:id", s.getAccount)
	auth.GET("/api/accounts", s.getAccounts)
	auth.GET("/api/accounts/del", s.getDeletedAccounts)
	auth.PATCH("/api/accounts/res/:id", s.restoreAccount)
	auth.DELETE("/api/accounts/:id", s.deleteAccount)

	// Transfer Routes
	auth.GET("/api/transfers/:id", s.getTransfers)
	auth.POST("/api/transfers", s.createTransfer)

	// User Routes
	auth.GET("api/users", s.getUser)
	auth.PATCH("api/users", s.updateUser)

	// Unauthenticated Routes
	router.POST("api/users/register", s.register)
	router.POST("api/users/login", s.loginUser)
	router.POST("api/users/renew", s.renewAccessToken)

	s.router = router
}

func (s *GinServer) setupSwagger() {
	s.router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
