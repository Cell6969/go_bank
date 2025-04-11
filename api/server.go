package api

import (
	db "github.com/Cell6969/go_bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP request for banking service
type Server struct {
	store  db.Store
	router *gin.Engine
}

// Create New Server instance
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// add routes to router
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts", server.listAccount)
	router.GET("/accounts/:id", server.getAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Start Server HTTP on specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Create error response
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
