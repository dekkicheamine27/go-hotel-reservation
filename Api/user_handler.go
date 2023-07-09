package Api

import "github.com/gofiber/fiber/v2"

func HandleGetUser(ctx *fiber.Ctx) error {
	return ctx.JSON("James")
}
