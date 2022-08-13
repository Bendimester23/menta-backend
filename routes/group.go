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

func HandleGroup_Delete(c *fiber.Ctx) error {
	if err := groupController.DeleteGroup(c.Locals(`id`).(string), c.Params(`id`)); err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.SendString("success")
}

func HandleGroup_ByID(c *fiber.Ctx) error {
	res, err := groupController.ById(c.Params(`id`), true)
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.JSON(res)
}

func HandleGroup_CreateExam(c *fiber.Ctx) error {
	var data models.CreateExam
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if err := data.Validate(); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := groupController.CreateExam(c.Params(`id`), data)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.JSON(res)
}

func HandleGroup_GetExams(c *fiber.Ctx) error {
	res, err := groupController.GetGroupExams(c.Params(`id`), c.Locals(`id`).(string))
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.JSON(res)
}

func HandleGroup_GetExamById(c *fiber.Ctx) error {
	res, err := groupController.GetExamById(c.Params(`id`), c.Locals(`id`).(string), c.Params(`eid`))
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.JSON(res)
}

func HandleGroup_DeleteExam(c *fiber.Ctx) error {
	err := groupController.DeleteExam(c.Locals(`id`).(string), c.Params(`id`), c.Params(`eid`))
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.SendString(`success`)
}
