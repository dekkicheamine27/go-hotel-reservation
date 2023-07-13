package main

import (
	"context"
	"github.com/godev/hotel-resevation/db"
	"github.com/godev/hotel-resevation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func hotelSeed(hotelName string, Location string, rating int) {
	hotel := types.Hotel{
		Name:     hotelName,
		Location: Location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	roms := []types.Room{
		{
			Size:  "small",
			Price: 99.99,
		},
		{
			Size:  "normal",
			Price: 122.9,
		},
		{
			Size:  "kingsize",
			Price: 150.9,
		},
	}

	insetHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range roms {
		room.HotelID = insetHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}

	}
}

func main() {
	hotelSeed("Safir", "Mostaganem", 3)
	hotelSeed("AZ", "Oran", 4)
	hotelSeed("sleep", "London", 1)

}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client, db.DBNAME)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
