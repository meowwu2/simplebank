package api

import (
	"errors"
	"fmt"
	"net/http"
	"simplebank/token"
	"strings"

	"github.com/gin-gonic/gin"
)

const(
	authorizationHeaderKey = "authorization"
	authorizationTypeBearer = "bearer"
)
func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc{
	return func (ctx *gin.Context)  {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader)==0{
			err:=errors.New("authorization headeer is not provided")
			ctx.AbortWithError(http.StatusUnauthorized,err)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields)<2{
			err:=errors.New("authorization headeer is not provided")
			ctx.AbortWithError(http.StatusUnauthorized,err)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer{
			err:=fmt.Errorf("unsupporteed authorization type %s",authorizationType)
			ctx.AbortWithError(http.StatusUnauthorized,err)
			return
		}

		accessToken := fields[1]
		payload,err := tokenMaker.VerifyToken(accessToken)
		if err!=nil{
			ctx.AbortWithError(http.StatusUnauthorized,err)
			return
		}
		ctx.Set(authorizationHeaderKey,payload)
		ctx.Next()
	}
}