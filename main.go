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

	var (
		hotelStore   = db.NewMongoHotelStore(client, DBNAME)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			User:    userStore,
			Hotel:   hotelStore,
			Room:    roomStore,
			Booking: bookingStore,
		}
		userHandler    = Api.NewUserHandler(userStore)
		hotelHandler   = Api.NewHotelHandler(store)
		authHandler    = Api.NewAuthHandler(userStore)
		roomHandler    = Api.NewRoomHandler(store)
		bookingHandler = Api.NewBookingHandler(store)

		app  = fiber.New(config)
		auth = app.Group("/api")

		apiV1 = app.Group("/api/v1", Api.JWTAuthentication(userStore))
		admin = apiV1.Group("/admin", Api.AdminAuth)
	)

	//auth api routes
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Api V1
	//user handlers
	apiV1.Get("/users", userHandler.HandleGetUsers)
	apiV1.Get("/users/:id", userHandler.HandleGetUser)
	apiV1.Post("/users", userHandler.HandlePostUser)
	apiV1.Delete("users/:id", userHandler.HandleDeleteUser)
	apiV1.Put("users/:id", userHandler.HandlePutUser)

	//hotel Handlers
	apiV1.Get("/hotels", hotelHandler.GetHotels)
	apiV1.Get("/hotels/:id", hotelHandler.GetHotel)
	apiV1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)

	//room Handlers
	apiV1.Get("/rooms", roomHandler.HandleGetRooms)
	apiV1.Post("/rooms/:id/book", roomHandler.HandleBookRoom)

	//booking handlers

	apiV1.Get("/bookings/:id", bookingHandler.HandleGetBooking)

	//admin handlers
	admin.Get("/bookings", bookingHandler.HandleGetBookings)

	err = app.Listen(":5000")
	if err != nil {
		log.Fatal(err)
	}
}
