package handler

import (

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler{
	return &UserHandler{usecase:usecase}
}

func (uh *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users , err := uh.usecase.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"Error":err.Error()})
	}

	return c.JSON(users)
}


func (uh *UserHandler) CreateUser(c *fiber.Ctx) error {
	var user *entity.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error":err.Error()})
	}

	if err := uh.usecase.CreateUser(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error":err.Error()})
	}

	return c.JSON(fiber.Map{"message":"user created successfully"})
}

