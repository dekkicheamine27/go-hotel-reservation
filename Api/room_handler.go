package Api

import (
	"context"
	"fmt"
	"github.com/godev/hotel-resevation/db"
	"github.com/godev/hotel-resevation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (b *BookRoomParams) validate() error {
	now := time.Now()
	if now.After(b.FromDate) || now.After(b.TillDate) {
		return fmt.Errorf("cannot book room in the past")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{store: store}
}

func (h *RoomHandler) HandleGetRooms(ctx *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(ctx.Context(), bson.M{})
	if err != nil {
		return err
	}
	return ctx.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(ctx *fiber.Ctx) error {
	var params BookRoomParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomID, err := primitive.ObjectIDFromHex(ctx.Params("id"))
	if err != nil {
		return err
	}
	user, ok := ctx.Context().Value("user").(*types.User)

	isAvailable, err := h.isRoomAvailable(ctx.Context(), roomID, params)
	if err != nil {
		return err
	}

	if !isAvailable {
		return ctx.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  "room already booked",
		})
	}
	if !ok {
		return err
	}
	booking := &types.Booking{
		RoomID:     roomID,
		UserID:     user.ID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	inserted, err := h.store.Booking.InsertBooking(ctx.Context(), booking)

	return ctx.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailable(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	fmt.Println(where)
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0
	return ok, err

}
