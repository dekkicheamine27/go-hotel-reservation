package db

const (
	DBNAME = "hotel-resevation"
	DURI   = "mongodb://localhost:27017"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
