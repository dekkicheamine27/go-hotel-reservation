package main

import (
	"context"
	"github.com/godev/hotel-resevation/Api"
	"github.com/godev/hotel-resevation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const dburi = "mongodb://localhost:27017"
const DBNAME = "hotel-resevation"

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := Api.NewUserHandler(db.NewMongoUserStore(client, DBNAME))

	app := fiber.New(config)
	apiV1 := app.Group("/api/v1")
	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Post("/users", userHandler.HandlePostUser)
	apiV1.Get("/users/:id", userHandler.HandleGetUser)
	apiV1.Delete("users/:id", userHandler.HandleDeleteUser)
	apiV1.Put("users/:id", userHandler.HandlePutUser)

	app.Listen(":5000")
}
