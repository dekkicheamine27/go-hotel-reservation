package db

import (
	"context"
	"fmt"
	"github.com/godev/hotel-resevation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingById(ctx context.Context, id string) (*types.Booking, error)
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{client: client, coll: client.Database(DBNAME).Collection("bookings")}
}

func (m *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	var bookings []*types.Booking
	curr, err := m.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	if err := curr.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (m *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := m.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(primitive.ObjectID)
	fmt.Println(booking.ID)
	return booking, nil
}

func (m *MongoBookingStore) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var booking types.Booking
	if err := m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}
