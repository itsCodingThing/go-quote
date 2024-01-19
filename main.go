package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/go-server/utils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	port, isPortExists := os.LookupEnv("PORT")
	if !isPortExists {
		port = "3000"
	}

	ctx := context.Background()
	fiberApp := fiber.New()

	client := utils.InitFirebase(ctx)
	defer client.Close()

	fiberApp.Get("/random", func(c *fiber.Ctx) error {
		response := utils.Response{}

		var offset int
		offset, err := strconv.Atoi(c.Query("offset"))
		if err != nil {
			response.StatusCode = 500
			response.Msg = "need offset"
			return response.CreateJSONResponse(c)
		}

		_, err = client.Collection("quotes").Offset(offset).Limit(5).Documents(ctx).GetAll()
		if err != nil {
			response.StatusCode = 500
			return response.CreateJSONResponse(c)
		}

		return response.CreateJSONResponse(c)
	})

	fmt.Println("server is running port:", port, "🚀")
	err := fiberApp.Listen(":" + port)
	if err != nil {
		fmt.Println("something went wrong")
	}
}
