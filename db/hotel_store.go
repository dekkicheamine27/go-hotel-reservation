package db

import (
	"context"
	"github.com/godev/hotel-resevation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	Update(ctx context.Context, filters bson.M, update bson.M) error
	GetHotel(ctx context.Context, m bson.M) ([]*types.Hotel, error)
	GetHotelById(ctx context.Context, id string) (*types.Hotel, error)
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbName string) *MongoHotelStore {
	return &MongoHotelStore{client: client, coll: client.Database(dbName).Collection("hotels")}
}

func (s *MongoHotelStore) GetHotel(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	curr, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err = curr.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil

}
func (m *MongoHotelStore) GetHotelById(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var hotel *types.Hotel
	if err := m.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel); err != nil {
		return nil, err
	}
	return hotel, nil
}

func (m *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := m.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filters bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filters, update)
	if err != nil {
		return err
	}
	return nil
}
