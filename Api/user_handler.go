package Api

import (
	"context"
	"errors"
	"github.com/godev/hotel-resevation/db"
	"github.com/godev/hotel-resevation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{userStore: userStore}
}

func (u *UserHandler) HandleDeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if err := u.userStore.DeleteUser(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.JSON(map[string]string{"delete": "true"})
}
func (u *UserHandler) HandlePutUser(ctx *fiber.Ctx) error {
	var params types.UpdateUserParams
	userId := ctx.Params("id")

	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	if err := u.userStore.UpdateUser(ctx.Context(), filter, params); err != nil {
		return err
	}
	return ctx.JSON(map[string]string{"msg": " update successful"})

}

func (u *UserHandler) HandlePostUser(ctx *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return ctx.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := u.userStore.InsertUser(ctx.Context(), user)
	if err != nil {
		return err
	}
	return ctx.JSON(insertedUser)
}

func (u *UserHandler) HandleGetUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	c := context.Background()
	user, err := u.userStore.GetUserById(c, id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ctx.JSON(map[string]string{"error": "not found"})
		}
		return err
	}
	return ctx.JSON(user)
}

func (u *UserHandler) HandleGetUsers(ctx *fiber.Ctx) error {
	var users []*types.User
	c := context.Background()
	users, err := u.userStore.GetUsers(c)
	if err != nil {
		return err
	}
	return ctx.JSON(users)
}
