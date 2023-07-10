package types

type User struct {
	ID        string `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string `bson:"firstName" json:"first_name"`
	LastName  string `bson:"lastName" json:"last_name"`
}
