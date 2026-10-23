package handler

import (
	"errors"
	"strconv"

	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	usecase usecase.UserUsecase
}

func NewUserHandler(usecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{usecase: usecase}
}

func (uh *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := uh.usecase.GetAllUsers()
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(users)
}

func (uh *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return badRequest(c, errors.New("invalid user id"))
	}

	user, err := uh.usecase.GetUserByID(uint(id))
	if err != nil {
		return notFoundOr500(c, err)
	}
	return c.JSON(user)
}

func (uh *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return badRequest(c, errors.New("invalid user id"))
	}

	if err := uh.usecase.DeleteUser(uint(id)); err != nil {
		return notFoundOr500(c, err)
	}
	return c.JSON(fiber.Map{"message": "successfully deleted"})
}
