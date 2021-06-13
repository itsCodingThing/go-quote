package main

import (
	"context"
	"log"
	"math/rand"
	"os"
	"time"

	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"github.com/imroc/req"
	"google.golang.org/api/option"
)

func main() {
	var port string

	value, isPortExists := os.LookupEnv("PORT")

	if !isPortExists {
		port = "1729"
	} else {
		port = value
	}
	ctx := context.Background()
	fiberApp := fiber.New()

	opt := option.WithCredentialsFile("./database.json")
	app, firebaseInitErr := firebase.NewApp(ctx, nil, opt)
	if firebaseInitErr != nil {
		log.Fatal(firebaseInitErr)
	}

	client, firstoreErr := app.Firestore(ctx)
	if firstoreErr != nil {
		log.Fatal(firstoreErr)
	}
	defer client.Close()

	fiberApp.Get("/", func(c *fiber.Ctx) error {
		rs, err := req.Get("https://node-quote.herokuapp.com/quote")

		if err != nil {
			return c.SendString("error")
		}

		return c.SendString(string(rs.Bytes()))
	})

	fiberApp.Get("/getIdList", func(c *fiber.Ctx) error {

		coll, err := client.Collection("quotes").Documents(ctx).GetAll()

		if err != nil {
			return c.SendString("something went wrong")
		}

		docSlice := make([]string, 0, 10)
		for _, doc := range coll {
			docSlice = append(docSlice, doc.Ref.ID)
		}

		docs := make(map[string]interface{})
		rand.Seed(time.Now().UnixNano())

		for i := 0; i < 20; i++ {
			docId := docSlice[rand.Intn(len(docSlice))]
			doc, _ := client.Collection("quotes").Doc(docId).Get(ctx)
			docs[docId] = doc.Data()
		}

		return c.JSON(map[string]interface{}{"dcos": docSlice})
	})

	fiberApp.Get("/random", func(c *fiber.Ctx) error {
		coll, err := client.Collection("quotes").Documents(ctx).GetAll()

		if err != nil {
			return c.SendString("something went wrong")
		}

		docSlice := make([]string, 0, 10)
		for _, doc := range coll {
			docSlice = append(docSlice, doc.Ref.ID)
		}

		// docs := make(map[string]interface{})
		// rand.Seed(time.Now().UnixNano())

		// for i := 0; i < 20; i++ {
		// 	docId := docSlice[rand.Intn(len(docSlice))]
		// 	doc, _ := client.Collection("quotes").Doc(docId).Get(ctx)
		// 	docs[docId] = doc.Data()
		// }

		docs := make([]map[string]interface{}, 0, 20)
		rand.Seed(time.Now().UnixNano())
		for i := 0; i < 20; i++ {
			docId := docSlice[rand.Intn(len(docSlice))]
			doc, _ := client.Collection("quotes").Doc(docId).Get(ctx)
			data := doc.Data()
			docs = append(docs, data)
		}

		return c.JSON(map[string]interface{}{"total": 20, "data": docs})
	})

	fiberApp.Get("/getjson", func(c *fiber.Ctx) error {
		return c.SendFile("./quote-database.json", true)
	})

	fiberApp.Listen(":" + port)
}
