package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	quote "example.com/server/utils"
	firebase "firebase.google.com/go"
	"github.com/imroc/req"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/api/option"
)

type Handler struct {
	ctx    context.Context
	client firestore.Client
}

func (h *Handler) GetRandomQuoteRoute(rq *req.Req) httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		res, err := rq.Get("https://node-quote.herokuapp.com/quote")

		if err != nil {
			log.Fatal(err)
		}

		rw.WriteHeader(http.StatusOK)
		rw.Write(res.Bytes())
	}

}

func (h *Handler) AddQuote() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		reqBody, _ := io.ReadAll(r.Body)
		body, _ := quote.UnmarshalQuote(reqBody)

		docRef, _, _ := h.client.Collection("qoutes-dev").Add(h.ctx, map[string]string{
			"title":   *body.Title,
			"content": *body.Content,
		})

		rw.Write([]byte(docRef.ID))
	}
}

func (h *Handler) GetQuoteById() httprouter.Handle {
	return func(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
		quoteId := p.ByName("id")

		docRef, _ := h.client.Collection("qoutes-dev").Doc(quoteId).Get(h.ctx)

		data := docRef.Data()

		j, _ := json.Marshal(data)

		rw.Write([]byte(j))
	}
}

func main() {
	port := os.Getenv("PORT")
	rq := req.New()
	ctx := context.Background()
	mux := httprouter.New()

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

	handler := Handler{
		ctx:    ctx,
		client: *client,
	}

	mux.GET("/", handler.GetRandomQuoteRoute(rq))
	mux.GET("/quote/:id", handler.GetQuoteById())
	mux.POST("/quote", handler.AddQuote())

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
