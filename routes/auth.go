package routes

import (
	"menta-backend/controller"
	"menta-backend/models"

	"github.com/gofiber/fiber/v2"
)

var authController = controller.AuthController{}

func HandleAuth_Register(c *fiber.Ctx) error {
	var data models.Register
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if err := data.Validate(); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := authController.RegisterUser(data)
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.JSON(fiber.Map{
		"id": res.ID,
		//FIXME: remove this
		"code": res.EmailCode,
	})
}

func HandleAuth_Verify(c *fiber.Ctx) error {
	err := authController.VerifyEmail(c.Params("code"), c.Params("id"))
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}
	return c.SendString("success")
}

func HandleAuth_Login(c *fiber.Ctx) error {
	var data models.Login
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if err := data.Validate(); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	res, err := authController.Login(data)
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	token, err := authController.CreateToken(res)
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	refreshToken, err := authController.CreateRefreshToken(res.ID)
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.JSON(fiber.Map{
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}

func HandleAuth_Refresh(c *fiber.Ctx) error {
	user, err := authController.UserById(c.Locals(`id`).(string))
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	token, err := authController.CreateToken(user)
	if err != nil {
		return c.Status(err.Code).SendString(err.Message)
	}

	return c.JSON(fiber.Map{
		"access_token": token,
	})
}
