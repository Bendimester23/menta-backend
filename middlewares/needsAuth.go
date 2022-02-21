package middlewares

import (
	"menta-backend/controller"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func NeedsAuth(c *fiber.Ctx) error {
	tokenStr := c.Get(`Authorization`)
	if len(tokenStr) == 0 {
		c.Locals(`auth`, false)
		return c.Status(401).SendString(`no token provided`)
	}
	tokenStr = tokenStr[7:]

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return controller.AuthTokenSecret, nil
	})

	if err != nil {
		c.Locals(`auth`, false)
		return c.Status(401).SendString(`wrong token provided`)
	}

	data := token.Claims.(jwt.MapClaims)

	if time.Since(time.UnixMilli(int64(data[`created`].(float64)*1000))).Seconds() >= 1800 {
		c.Locals(`auth`, false)
		return c.Status(401).SendString(`expired token`)
	}

	c.Locals(`id`, data[`id`])
	c.Locals(`username`, data[`username`])
	c.Locals(`created`, data[`created`])

	c.Locals(`auth`, true)
	return c.Next()
}
