package routes

import (
	"menta-backend/controller"

	"github.com/gofiber/fiber/v2"
)

var userController = controller.UserController{}

func HandleUser_Me(c *fiber.Ctx) error {
	res, err := userController.GetUserById(c.Locals(`id`).(string))
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.JSON(res)
}

func HandleUser_Groups(c *fiber.Ctx) error {
	res, err := userController.GetGroups(c.Locals(`id`).(string))
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.JSON(res)
}

func HandleUser_JoinGroup(c *fiber.Ctx) error {
	err := userController.JoinGroup(c.Locals(`id`).(string), c.Query(`code`))
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.SendString(`success`)
}
