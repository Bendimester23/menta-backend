package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"menta-backend/db"
	"menta-backend/mail"
	"menta-backend/middlewares"
	"menta-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	db.Connect()
	defer db.Disconnect()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	app.Use(func(c *fiber.Ctx) error {
		// start timer
		start := time.Now().UnixMilli()
		// next routes
		err := c.Next()
		// stop timer
		stop := time.Now().UnixMilli()
		// Do something with response
		c.Append("Server-timing", fmt.Sprintf("app;dur=%vms", stop-start))
		return err
	})

	app.Use(logger.New())

	mail.InitDialer()

	initRoutes(app.Group("/api"))

	if err := app.Listen(fmt.Sprintf(":%s", getPort())); err != nil {
		log.Fatalln(err)
	}
}

func getPort() string {
	if os.Getenv("PORT") == "" {
		return "8080"
	}
	return os.Getenv("PORT")
}

func initRoutes(r fiber.Router) {
	auth := r.Group("/auth")
	auth.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: time.Minute,
	}))
	auth.Put("/register", routes.HandleAuth_Register)
	auth.Post("/login", routes.HandleAuth_Login)
	auth.Get("/verify/:id/:code", routes.HandleAuth_Verify)
	auth.Use(middlewares.NeedsRefreshToken)
	auth.Post("/refresh", routes.HandleAuth_Refresh)

	user := r.Group("/user")
	user.Use(middlewares.NeedsAuth)
	user.Get("/me", routes.HandleUser_Me)
	user.Get(`/groups`, routes.HandleUser_Groups)
	user.Get(`/join-group`, routes.HandleUser_JoinGroup)

	avatars := r.Group("/avatar")
	avatars.Use(filesystem.New(filesystem.Config{
		Root:         http.Dir(`./avatars`),
		Browse:       false,
		Index:        `default.svg`,
		NotFoundFile: `default.svg`,
	}))

	group := r.Group("/group")
	group.Use(middlewares.NeedsAuth)
	group.Put("/create", routes.HandleGroup_Create)
	group.Get(`/:id/`, routes.HandleGroup_ByID)
}
