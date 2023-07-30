package Api

import (
	"fmt"
	"github.com/godev/hotel-resevation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		fmt.Println("--JWT authentication")
		token, ok := ctx.GetReqHeaders()["X-Api-Token"]
		if !ok {
			return ErrorUnauthorized()
		}

		claims, err := parseJWTToken(token)

		if err != nil {
			return err
		}

		// Check token expiration
		userID := claims["id"].(string)
		user, err := userStore.GetUserById(ctx.Context(), userID)
		if err != nil {
			return ErrorUnauthorized()
		}
		ctx.Context().SetUserValue("user", user)

		return ctx.Next()
	}
}

func parseJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signature method", token.Header["alg"])
			return nil, ErrorUnauthorized()
		}

		secret := os.Getenv("JWT_SECRET")
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, ErrorUnauthorized()
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		fmt.Println(claims["foo"], claims["nbf"])
	}
	return claims, nil
}
