package main

import (
	h "cards/internal/api/v1/handlers"
	"cards/internal/database"
	"cards/internal/repository"
	"cards/internal/service"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

func main() {
	// Initialize dependencies
	db, err := database.NewDatabaseConnection(false)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	cardsRepository := repository.NewCardsRepository(db)
	cardsService := service.NewCardsService(service.CardsService{Repository: cardsRepository})

	ch := h.NewCardsHandler(nil, cardsService)
	ch.Handlers()

	// Initialize the server
	log.Print("Listening on port 3000 🚀")
	err = http.ListenAndServe(":3000", nil)
	log.Fatal(err)
}

func mainAWS() {
	cardsService := service.NewCardsService(service.CardsService{Repository: nil})
	_ = h.NewCardsHandler(nil, cardsService)
	log.Print("Starting lambda 🚀")
	lambda.Start(httpadapter.New(http.DefaultServeMux).ProxyWithContext)
}