package middleware

import (
	"github/mgabrielo/React-Golang-Blog/util"

	"github.com/gofiber/fiber/v2"
)

func IsAuthenticate(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	if _, err := util.ParseJWT(cookie); err != nil {
		c.Status(401)
		return c.JSON(fiber.Map{"message": "unauthorised access"})
	}
	return c.Next()
}
