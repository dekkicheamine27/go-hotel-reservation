package Api

import (
	"context"
	"github.com/godev/hotel-resevation/db"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}
}

func (u *UserHandler) HandleGetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	c := context.Background()
	user, err := u.userStore.GetUserById(c, id)
	if err != nil {
		return err
	}
	return ctx.JSON(user)
}

func (u *UserHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	return ctx.JSON("James")
}
