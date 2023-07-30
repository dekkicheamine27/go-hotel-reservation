package Api

import (
	"context"
	"github.com/godev/hotel-resevation/db"
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	dbName := os.Getenv("MONGO_DB_NAME")
	if err := tdb.client.Database(dbName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}
	mongodbUri := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongodbUri))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,
		Store: &db.Store{
			Hotel:   hotelStore,
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}
