package response

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/alireza-hmd/c2/pkg/encrypt/aes"
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

func EncryptedErrorResponse(c *gin.Context, code int, msg string, token string) {
	res := Response{
		Success: false,
		Message: msg,
	}
	data, err := json.Marshal(&res)
	if err != nil {
		log.Println(err)
		return
	}
	key := aes.StrToByte(token)
	data, err = aes.Encrypt(data, key)

	c.String(code, "%s", string(data))
}

func OkResponse(c *gin.Context, msg string, body interface{}) {
	res := Response{
		Message: msg,
		Success: true,
		Data:    body,
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

func EncryptedOkResponse(c *gin.Context, msg string, token string, body interface{}) {
	res := Response{
		Success: true,
		Message: msg,
		Data:    body,
	}
	data, err := json.Marshal(&res)
	if err != nil {
		log.Println(err)
		return
	}
	key := aes.StrToByte(token)
	data, err = aes.Encrypt(data, key)

	c.String(200, "%s", string(data))
}
