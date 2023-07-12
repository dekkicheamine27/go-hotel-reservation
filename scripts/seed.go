package main

import (
	"context"
	"fmt"
	"github.com/godev/hotel-resevation/db"
	"github.com/godev/hotel-resevation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME)

	hotel := types.Hotel{
		Name:     "Safir",
		Location: "Mostaganem",
	}

	roms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 9.99,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 199.9,
		},
		{
			Type:      types.DoubleRoomType,
			BasePrice: 120.9,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 220.9,
		},
	}

	insetHotel, err := hotelStore.InsertHotel(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insetHotel)

	for _, room := range roms {
		room.HotelID = insetHotel.ID
		insertRoom, err := roomStore.InsertRoom(context.Background(), &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertRoom)

	}

}
