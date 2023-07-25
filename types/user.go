package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p *UpdateUserParams) ToBison() bson.M {
	m := bson.M{}
	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}
	return m
}

func (params CreateUserParams) Validate() map[string]string {
	var errors = map[string]string{}
	if len(params.FirstName) < 3 {
		errors["firstName"] = "FirstName must be at least 3 characters"
	}

	if len(params.LastName) < 3 {
		errors["lastName"] = "LastName must be at least 3 characters"
	}
	if len(params.Password) < 8 {
		errors["password"] = "Password must be at least 8 characters"
	}

	if !isEmailValid(params.Email) {
		errors["email"] = "Email is invalid"
	}

	return errors

}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	crypt, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(crypt),
	}, nil
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName         string             `bson:"firstName" json:"first_name"`
	LastName          string             `bson:"lastName" json:"last_name"`
	Email             string             `bson:"email" json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
	IsAdmin           bool               `bson:"isAdmin" json:"isAdmin"`
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}
