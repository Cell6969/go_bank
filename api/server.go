package api

import (
	"fmt"

	db "github.com/Cell6969/go_bank/db/sqlc"
	"github.com/Cell6969/go_bank/token"
	"github.com/Cell6969/go_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP request for banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// Create New Server instance
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

// Setup Router
func (server *Server) setupRouter() {
	router := gin.Default()

	// Add routes to router
	// Account Route
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccount)
	router.GET("/accounts/:id", server.getAccount)

	// Transfer Route
	router.POST("/transfers", server.createTransfer)

	// User Route
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	server.router = router
}

// Start Server HTTP on specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Create error response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
