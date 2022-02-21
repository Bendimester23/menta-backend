package routes

import (
	"menta-backend/controller"
	"menta-backend/models"

	"github.com/gofiber/fiber/v2"
)

var groupController = controller.GroupController{}

func HandleGroup_Create(c *fiber.Ctx) error {
	var data models.CreateGroup
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if err := data.Validate(); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := groupController.Create(c.Locals(`id`).(string), data)
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.JSON(res)
}

func HandleGroup_ByID(c *fiber.Ctx) error {
	res, err := groupController.ById(c.Params(`id`), true)
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.JSON(res)
}
