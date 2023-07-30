package Api

import (
	"fmt"
	"github.com/godev/hotel-resevation/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthUser(ctx *fiber.Ctx) (*types.User, error) {
	user, ok := ctx.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("unauorized")
	}
	return user, nil
}
