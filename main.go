package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"menta-backend/db"
	"menta-backend/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	app.Use(logger.New())

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
	//TODO: add middleware
	auth.Post("/refresh", routes.HandleAuth_Refresh)
	routes.InitAuthWs(auth)

	user := r.Group("/user")
	user.Get("/me", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"user": fiber.Map{
				"id":        "ckz5vy1dy00021st5f2vqzkxo",
				"createdAt": "2022-02-02T18:30:50.758Z",
				"updatedAt": "2022-02-02T18:31:01.161Z",
				"username":  "teszt69",
				"email":     "cucc@gmail.com",
				"isTeacher": false,
			},
		})
	})
}
