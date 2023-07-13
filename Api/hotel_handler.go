package Api

import (
	"context"
	"github.com/godev/hotel-resevation/db"
	"github.com/godev/hotel-resevation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{store: store}
}

func (h *HotelHandler) GetHotels(ctx *fiber.Ctx) error {

	var hotels []*types.Hotel
	c := context.Background()
	hotels, err := h.store.Hotel.GetHotel(c, bson.M{})
	if err != nil {
		return err
	}
	return ctx.JSON(hotels)
}

func (h *HotelHandler) HandleGetRooms(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	c := context.Background()
	rooms, err := h.store.Room.GetRooms(c, bson.M{"hotelID": oid})
	if err != nil {
		return err
	}

	return ctx.JSON(rooms)

}

func (h *HotelHandler) GetHotel(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	c := context.Background()
	hotel, err := h.store.Hotel.GetHotelById(c, id)
	if err != nil {
		return err
	}
	return ctx.JSON(hotel)

}
