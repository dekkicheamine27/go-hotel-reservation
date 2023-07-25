package Api

import (
	"fmt"
	"github.com/godev/hotel-resevation/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(ctx *fiber.Ctx) error {
	user, ok := ctx.Context().UserValue("user").(*types.User)
	if !ok || !user.IsAdmin {
		return fmt.Errorf("not authorized")
	}

	return ctx.Next()

}
