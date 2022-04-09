package middlewares

import (
	"context"
	"errors"
	"menta-backend/controller"
	"menta-backend/db"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func NeedsRefreshToken(c *fiber.Ctx) error {
	parts := strings.Split(string(c.Body()), `"`)
	if len(parts) != 5 {
		c.Locals(`auth`, false)
		return c.Status(400).SendString(`wrong token provided`)
	}
	tokenStr := parts[3]
	if len(tokenStr) == 0 {
		c.Locals(`auth`, false)
		return c.Status(401).SendString(`no token provided`)
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return controller.RefreshTokenSecret, nil
	})

	if err != nil {
		c.Locals(`auth`, false)
		return c.Status(401).SendString(`wrong token provided`)
	}

	data := token.Claims.(jwt.MapClaims)

	if data[`type`].(string) != `refresh` {
		c.Locals(`auth`, false)
		return c.Status(401).SendString(`wrong token provided`)
	}

	if time.Since(time.UnixMilli(int64(data[`created`].(float64)*1000))).Seconds() >= 60*60*24*30 {
		c.Locals(`auth`, false)
		return c.Status(401).SendString(`expired token`)
	}

	res, err := db.DB.RefreshToken.FindMany(
		db.RefreshToken.ID.Equals(data[`refresh_id`].(string)),
	).With(
		db.RefreshToken.Owner.Fetch(),
	).Exec(context.Background())

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			c.Locals(`auth`, false)
			return c.Status(401).SendString(`invalid token provided`)
		}
		c.Locals(`auth`, false)
		return c.Status(500).SendString(`db error`)
	}

	if len(res) == 0 {
		c.Locals(`auth`, false)
		return c.Status(401).SendString(`invalid token provided`)
	}

	c.Locals(`auth`, true)
	c.Locals(`id`, res[0].Owner().ID)
	c.Locals(`username`, res[0].Owner().Username)

	return c.Next()
}
