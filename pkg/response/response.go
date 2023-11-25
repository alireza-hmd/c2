package response

import (
	"errors"
)

var ErrNotFound = errors.New("not found\n")
var ErrInvalidRequest = errors.New("invalid request body\n")

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
	Error   Error       `json:"error"`
	Data    interface{} `json:"data"`
}

type Error struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
}
