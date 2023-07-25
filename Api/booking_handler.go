package Api

import (
	"github.com/godev/hotel-resevation/db"
	"github.com/godev/hotel-resevation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{store: store}
}

func (h *BookingHandler) HandleGetBookings(ctx *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(ctx.Context(), bson.M{})
	if err != nil {
		return err
	}
	return ctx.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingById(ctx.Context(), id)
	if err != nil {
		return err
	}

	user, ok := ctx.Context().UserValue("user").(*types.User)
	if !ok {
		return err
	}

	if user.ID != booking.UserID {
		return ctx.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}

	return ctx.JSON(booking)

}
