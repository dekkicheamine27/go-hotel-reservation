package main

import (
	"github.com/godev/hotel-resevation/types"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get("/foo", handleFoo)
	apiV1 := app.Group("/api/v1")
	apiV1.Get("/user", handleUser)
	app.Listen(":5000")
}

func handleUser(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "Dekkiche",
		LastName:  "Amine",
	}
	return c.JSON(u)
}

func handleFoo(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"msg": "is working fine!"})
}
