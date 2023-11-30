package response

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var ErrNotFound = errors.New("not found\n")
var ErrInvalidRequest = errors.New("invalid request body\n")

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ErrorResponse(c *gin.Context, code int, msg string) {
	res := Response{
		Success: false,
		Message: msg,
	}
	c.IndentedJSON(code, res)
}

func OkResponse(c *gin.Context, msg string) {
	res := Response{
		Message: msg,
		Success: true,
	}
	c.IndentedJSON(200, res)
}

func OkResponseWithData(c *gin.Context, msg string, data interface{}) {
	res := Response{
		Message: msg,
		Success: true,
		Data:    data,
	}
	c.IndentedJSON(200, res)
}
