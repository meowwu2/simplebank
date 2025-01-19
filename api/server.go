package api

import (
	"fmt"
	db "simplebank/db/sqlc"
	"simplebank/token"
	"simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP requests for our banking service
type Server struct{
	config util.Config
	store db.Store
	tokenMaker token.Maker
	router *gin.Engine
}

// NewServer creates a new HTTP server and set up routing
func NewServer(config util.Config,store db.Store)(*Server,error){
	tokenMaker,err:=token.NewPasetoMaker(config.TOKEN_SYMMETRIC_KEY)
	if err!=nil{
		return nil,fmt.Errorf("cannot create token maker: %w",err)
	}
	server := &Server{
		config: config ,
		tokenMaker: tokenMaker,
		store: store,
	}
	server.setupServer()

	if v,ok:=binding.Validator.Engine().(*validator.Validate);ok{
		v.RegisterValidation("currency",validCurrency)
	}
	
	return server,nil
}

func (server *Server)setupServer()  {
	router := gin.Default()
	router.POST("/user",server.CreateUser)
	router.POST("/user/login",server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))
	authRoutes.POST("/account",server.CreateAccount)
	authRoutes.GET("/account/:id",server.GetAccount)
	authRoutes.GET("/account",server.ListAccount)
	authRoutes.POST("/transfer",server.CreateTransfer)

	server.router=router
}
//Start runs the HTTP server on a specific address
func(server *Server)Start(address string)error{
	return server.router.Run(address)
}
func errorResponse(err error)gin.H{
	return gin.H{"error":err.Error()}
}