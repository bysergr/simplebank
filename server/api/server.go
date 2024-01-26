package api

import (
	"errors"
	"net/http"

	db "github.com/bysergr/simple-bank/db/sqlc"
	"github.com/bysergr/simple-bank/token"
	"github.com/bysergr/simple-bank/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     utils.Config
}

var ErrConvertPayload error = errors.New("failed in convert authorization payload")

func NewServer(store db.Store, tokenMaker token.Maker, config utils.Config) *Server {
	server := &Server{store: store, tokenMaker: tokenMaker, config: config}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)
	router.POST("/token/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

func errorPayload(ctx *gin.Context) (authPayload *token.Payload, valid bool) {
	authPayload, ok := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, errorResponse(ErrConvertPayload))
		return authPayload, false
	}

	return authPayload, true
}
