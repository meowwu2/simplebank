package api

import (
	"database/sql"
	"errors"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type CreateAccountRequest struct{
	Currency string `json:"currency" binding:"required,currency"`
} 

func (server *Server) CreateAccount(ctx *gin.Context){
	var req CreateAccountRequest
	if err := ctx.ShouldBindJSON(&req);err!=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return 
	}

	authPayload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	arg :=db.CreateAccountParams{
		Owner:authPayload.Username,
		Currency: req.Currency,
		Balance: 0,
	}
	account,err:=server.store.CreateAccount(ctx,arg)
	if err!=nil{
		if pqErr,ok:=err.(*pq.Error);ok{
			if pqErr.Code.Name() == "foreign_key_violation" || pqErr.Code.Name() == "unique_violation" {
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,account)
}

type getAccountRequest struct{
	ID int64 `uri:"id" binding:"required,min=1" `
}
func (server *Server) GetAccount(ctx *gin.Context){
	var req getAccountRequest
	if err := ctx.ShouldBindUri(&req);err!=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return 
	}

	account,err:=server.store.GetAccount(ctx,req.ID)
	if err!=nil{
		if err==sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	if account.Owner != authPayload.Username{
		err:=errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,account)

}

type listAccountRequest struct{
	PageID   int32 `form:"page_id" binding:"required,min=1" `
	PageSize int32 `form:"page_size" binding:"required,min=1,max=10" `

}
func (server *Server) ListAccount(ctx *gin.Context){
	var req listAccountRequest
	if err := ctx.ShouldBindQuery(&req);err!=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return 
	}
	authPayload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	arg:=db.ListAccountsParams{
		Owner: authPayload.Username,
		Limit: req.PageSize,
		Offset: (req.PageID-1)*req.PageSize,
	}
	accounts,err:=server.store.ListAccounts(ctx,arg)
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,accounts)

}
