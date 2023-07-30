package Api

import (
	"encoding/json"
	"fmt"
	"github.com/godev/hotel-resevation/db/fixtures"
	"github.com/godev/hotel-resevation/types"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetBooking(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)
	var (
		nonAuthUser    = fixtures.AddUser(db.Store, "dekk", "moh", false)
		user           = fixtures.AddUser(db.Store, "dekkiche", "amine", false)
		hotel          = fixtures.AddHotel(db.Store, "zohor", "Mostaganem", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		jwt            = app.Group("/", JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)
	jwt.Get("/:id", bookingHandler.HandleGetBooking)
	res := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	res.Header.Add("X-Api-Token", CreateToken(user))

	resp, err := app.Test(res)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 status code, got %d", resp.StatusCode)
	}

	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}

	if booking.ID != bookingResp.ID {
		t.Fatalf("expected %s got %s", booking.ID, bookingResp.ID)
	}
	if bookingResp.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, bookingResp.UserID)
	}

	// test non-auth user cannot access the booking
	res = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	res.Header.Add("X-Api-Token", CreateToken(nonAuthUser))
	resp, err = app.Test(res)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		adminUser      = fixtures.AddUser(db.Store, "admin", "admin", true)
		user           = fixtures.AddUser(db.Store, "dekkiche", "amine", false)
		hotel          = fixtures.AddHotel(db.Store, "zohor", "Mostaganem", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 4.4, hotel.ID)
		from           = time.Now()
		till           = from.AddDate(0, 0, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, from, till)
		app            = fiber.New()
		admin          = app.Group("/", JWTAuthentication(db.User), AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)
	admin.Get("/", bookingHandler.HandleGetBookings)
	res := httptest.NewRequest("GET", "/", nil)
	res.Header.Add("X-Api-Token", CreateToken(adminUser))

	resp, err := app.Test(res)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 status code, got %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Fatalf("expected %s got %s", booking.ID, have.ID)
	}
	if have.UserID != booking.UserID {
		t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
	}

	// test non-admin cannot access the bookings
	admin.Get("/", bookingHandler.HandleGetBookings)
	res = httptest.NewRequest("GET", "/", nil)
	res.Header.Add("X-Api-Token", CreateToken(user))

	resp, err = app.Test(res)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected not 200 status code, got %d", resp.StatusCode)
	}
}
