package handler

import (
	"strconv"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
)

type CharHandler struct {
	usecase usecase.CharUsecase
}

func NewCharHandler(usecase usecase.CharUsecase) *CharHandler {
	return &CharHandler{usecase}
}

func (h *CharHandler) GetAllCharsByProjectID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid project ID",
			"error":   err.Error(),
		})
	}
	chars, err := h.usecase.GetAllCharsByProjectID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting characters by project ID",
			"error":   err.Error(),
		})
	}

	return c.JSON(chars)

}

func (h *CharHandler) CreateChar(c *fiber.Ctx) error {
	var char entity.Char

	if err := c.BodyParser(&char); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
			"error":   err.Error(),
		})
	}

	err := h.usecase.CreateChar(&char)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating character",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(char)
}

func (h *CharHandler) UpdateChar(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid character ID",
			"error":   err.Error(),
		})
	}

	var char entity.Char

	if err := c.BodyParser(&char); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Error parsing request body",
			"error":   err.Error(),
		})
	}

	err = h.usecase.UpdateChar(&char, uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating character",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{"message": "Char updated"})
}

func (h *CharHandler) DeleteChar(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid character ID",
			"error":   err.Error(),
		})
	}

	err = h.usecase.DeleteChar(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting character",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{"message": "Char deleted"})
}
