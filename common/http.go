package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 该函数返回一个gin.H，gin.H是一个map，存储着键值对，将要返回给请求者
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func GetGinHandler[REQ any, RES any](logic func(req *REQ) (*RES, error)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var req REQ
		if err := ctx.ShouldBindJSON(&req); err != nil {
			//证明请求对于该结构体并不有效
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}
		res, err := logic(&req)
		if err != nil {
			ctx.JSON(http.StatusPreconditionFailed, errorResponse(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}
