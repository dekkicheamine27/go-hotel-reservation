package Api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/godev/hotel-resevation/db"
	"github.com/godev/hotel-resevation/types"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"reflect"

	"net/http/httptest"
	"testing"
)

func insertUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "dekkiche",
		LastName:  "rayhane",
		Email:     "dekkichos@gmail.com",
		Password:  "123456",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}

	return user
}

func TestJWTAuthentication(t *testing.T) {
	tdb := setUp(t)
	defer tdb.teardown(t)
	insertedUser := insertUser(t, tdb.UserStore)

	app := fiber.New()
	authHAndler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHAndler.HandleAuthenticate)

	params := &AuthParams{
		Email:    "dekkichos@gmail.com",
		Password: "123456",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected %v, got %v", http.StatusOK, resp.Status)
	}

	var authResp AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		t.Error(err)
	}

	if authResp.Token == "" {
		t.Fatalf("No token")
	}

	// Set the encrypted password to an empty string, because we do NOT return that in any
	insertedUser.EncryptedPassword = ""
	if !reflect.DeepEqual(authResp.User, insertedUser) {
		t.Fatal("expected user is inserted user")
	}

}

func TestJWTAuthenticationWithWrongPassword(t *testing.T) {
	tdb := setUp(t)
	defer tdb.teardown(t)
	insertUser(t, tdb.UserStore)

	app := fiber.New()
	authHAndler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHAndler.HandleAuthenticate)

	params := &AuthParams{
		Email:    "dekkichos@gmail.com",
		Password: "12345678",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected %v, got %v", http.StatusBadRequest, resp.Status)
	}

	var genResp *genericResp
	err = json.NewDecoder(resp.Body).Decode(&genResp)
	if err != nil {
		t.Error(err)
	}

	if genResp.Type != "error" {
		t.Fatalf("expected gen response type to be error but got %s", genResp.Type)

	}

	if genResp.Msg != "invalid credentials" {
		t.Fatalf("expected gen response msg to be <invalid credentials> but got %s", genResp.Msg)
	}

}
