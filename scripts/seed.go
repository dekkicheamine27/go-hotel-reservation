package main

import (
	"context"
	"fmt"
	"github.com/godev/hotel-resevation/Api"
	"github.com/godev/hotel-resevation/db"
	"github.com/godev/hotel-resevation/db/fixtures"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	client *mongo.Client
	ctx    = context.Background()
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var err error
	mongodbUrl := os.Getenv("MONGO_DB_URL")
	dbName := os.Getenv("MONGO_DB_NAME")
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodbUrl))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(dbName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(store, "dekkiche", "amine", false)
	fmt.Println("amine ->", Api.CreateToken(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin ->", Api.CreateToken(admin))
	hotel := fixtures.AddHotel(store, "AZ", "Mostaganem", 5, nil)
	room := fixtures.AddRoom(store, "large", true, 88.44, hotel.ID)
	booking := fixtures.AddBooking(store, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 0, 5))
	fmt.Println("booking ->", booking.ID)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}
