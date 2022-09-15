package routes

import (
	"menta-backend/controller"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

var chatController = new(controller.ChatController)

func HandleChat_GetMessages(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit > 100 {
		return c.Status(400).SendString("Parameter \"limit\" is not a number or it is higher than 100.")
	}

	res, errRes := chatController.GetMessages(c.Params(":id"), limit)

	if errRes != nil {
		return c.Status(errRes.Code).SendString(errRes.Message)
	}

	return c.JSON(res)
}

func HandleChat_SendMessage(c *fiber.Ctx) error {
	return c.SendString("igen")
}
