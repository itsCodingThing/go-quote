package utils

import (
	"github.com/gofiber/fiber/v2"
)

// Response struct represents a Response object
type Response struct {
	Msg        string      `json:"msg"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

func (response *Response) CreateJSONResponse(ctx *fiber.Ctx) error {
	if response.Data == nil {
		response.Data = ""
	}

	if response.Msg == "" {
		response.Msg = "success"
	}

	if response.StatusCode == 0 {
		response.StatusCode = 200
	}

	return ctx.Status(response.StatusCode).JSON(response)
}
