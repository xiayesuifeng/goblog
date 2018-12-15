package core

import "github.com/gin-gonic/gin"

const (
	ResultSuccessCode             = 200
	ResultFailCode                = 400
	ResultUnauthorizedCode        = 401
	ResultNotFoundCode            = 404
	ResultInternalServerErrorCode = 500
)

func Result(code int, message string) gin.H {
	return gin.H{
		"code":    code,
		"message": message,
	}
}

func SuccessResult() gin.H {
	return gin.H{
		"code": ResultSuccessCode,
	}
}

func SuccessDataResult(name string, data interface{}) gin.H {
	return gin.H{
		"code": ResultSuccessCode,
		name:   data,
	}
}

func FailResult(message string) gin.H {
	return gin.H{
		"code":    ResultFailCode,
		"message": message,
	}
}
