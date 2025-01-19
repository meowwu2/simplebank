package api

import (
	"database/sql"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Username       string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	FullName       string `json:"full_name" binding:"required"`
	Email          string `json:"email" binding:"required,email" `
}
type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func NewUserResponse(user db.User) userResponse {
	return userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email: user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) CreateUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	HashedPassword,err := util.HashPassword(req.Password)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))

	}
	arg := db.CreateUserParams{
		Username: req.Username,
		FullName: req.FullName,
		HashedPassword:HashedPassword ,
		Email: req.Email,
	}
	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code.Name() == "foreign_key_violation" || pqErr.Code.Name() == "unique_violation" {
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp:= NewUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type loginUserRequest struct {
	Username       string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken 	string 			`json:"access_token"`
	User 			userResponse 	`json:"user"`
}

func (server *Server)loginUser(ctx *gin.Context)  {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	usr,err := server.store.GetUser(ctx,req.Username)
	if err !=nil{
		if err == sql.ErrNoRows{
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.Password,usr.HashedPassword)
	if err != nil{
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}
	accessToken,err:=server.tokenMaker.CreateToken(req.Username,server.config.ACCESS_TOKEN_DURATION)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	rsp:=loginUserResponse{
		AccessToken: accessToken,
		User: NewUserResponse(usr),
	}
	ctx.JSON(http.StatusOK,rsp)
}

