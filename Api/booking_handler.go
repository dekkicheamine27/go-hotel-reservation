package Api

import (
	"fmt"
	"github.com/godev/hotel-resevation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
		return ErrNotResourceNotFound("bookings")
	}
	return ctx.JSON(bookings)
}

func (h *BookingHandler) HandleGetBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingById(ctx.Context(), id)
	if err != nil {
		return ErrNotResourceNotFound("booking")
	}

	user, err := getAuthUser(ctx)
	if err != nil {
		return ErrorUnauthorized()
	}

	if user.ID != booking.UserID {
		return ErrorUnauthorized()
	}

	return ctx.JSON(booking)

}

func (h *BookingHandler) HandleCancelBooking(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	booking, err := h.store.Booking.GetBookingById(ctx.Context(), id)
	if err != nil {
		return err
	}

	user, err := getAuthUser(ctx)
	if err != nil {
		return ErrorUnauthorized()
	}

	if user.ID != booking.UserID {
		return ErrorUnauthorized()
	}

	if err := h.store.Booking.UpdateBooking(ctx.Context(), ctx.Params("id"), bson.M{"cancelled": true}); err != nil {
		return fmt.Errorf("update book not working")
	}

	return ctx.JSON(genericResp{Msg: "updated", Type: "msg"})

}
