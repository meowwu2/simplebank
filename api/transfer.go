package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	db "simplebank/db/sqlc"
	"simplebank/token"

	"github.com/gin-gonic/gin"
)

type transferRequest struct{
	From_Account_ID     int64 `json:"from_account_id" binding:"required,min=1"` 
	To_Account_ID    	int64 `json:"to_account_id" binding:"required,min=1"` 
	Amount				int64  `json:"amount" binding:"required,gt=0"`
	Currency 			string `json:"currency" binding:"required,currency"`
} 

func (server *Server) CreateTransfer(ctx *gin.Context){
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req);err!=nil{
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return 
	}
	fromAccount,valid :=server.validAccountCurrency(ctx,req.From_Account_ID,req.Currency)
	authPayload := ctx.MustGet(authorizationHeaderKey).(*token.Payload)
	if fromAccount.Owner != authPayload.Username{
		err:=errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized,errorResponse(err))
		return
	}
	if !valid{
		return
	}
	_,valid =server.validAccountCurrency(ctx,req.To_Account_ID,req.Currency)

	if !valid{
		return
	}


	arg :=db.TransferTxParams{
		FromAccountID:	req.From_Account_ID,
		ToAccountID: 	req.To_Account_ID,
		Amount: 		req.Amount,
	}
	result,err:=server.store.TransferTx(ctx,arg)
	
	if err!=nil{
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK,result)
}

func (server *Server)validAccountCurrency(ctx *gin.Context,accounID int64,currency string)(db.Account,bool){
	account,err:=server.store.GetAccount(ctx,accounID)
	if err!=nil{
		if err==sql.ErrNoRows{
			ctx.JSON(http.StatusNotFound,errorResponse(err))
			return account,false
		}
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return account,false
	}

	if account.Currency!=currency{
		err := fmt.Errorf("account[%d] currency mismatch:%s vs %s",accounID,account.Currency,currency)
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return account,false
	}
	return account,true
}