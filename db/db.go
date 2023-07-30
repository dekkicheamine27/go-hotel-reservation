package db

import "os"

var (
	DBNAME string
	DURI   string
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

func init() {
	DBNAME = os.Getenv("MONGO_DB_NAME")
	DURI = os.Getenv("MONGO_DB_URL")
}
