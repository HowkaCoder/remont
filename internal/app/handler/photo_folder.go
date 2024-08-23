package handler

import (
	"strconv"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
)

type PhotoFolderHandler struct {
	usecase usecase.PhotoFolderUseCase
}

func NewPhotoFolderHandler(usecase usecase.PhotoFolderUseCase) *PhotoFolderHandler {
	return &PhotoFolderHandler{usecase: usecase}
}
func (h *PhotoFolderHandler) GetAllPhotoFolders(c *fiber.Ctx) error {
	folders, err := h.usecase.GetAllPhotoFolders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting photo folders",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": folders,
	})
}

func (h *PhotoFolderHandler) GetPhotoFolderByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid folder ID",
			"error":   err.Error(),
		})
	}

	folder, err := h.usecase.GetPhotoFolderByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting folder by ID",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": folder,
	})
}

func (h *PhotoFolderHandler) GetPhotoFoldersByProjectID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid project ID",
			"error":   err.Error(),
		})
	}

	folders, err := h.usecase.GetPhotoFoldersByProjectID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting folders by project ID",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": folders,
	})
}

func (h *PhotoFolderHandler) CreatePhotoFolder(c *fiber.Ctx) error {
	var folder entity.PhotoFolder
	if err := c.BodyParser(&folder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	if err := h.usecase.CreatePhotoFolder(&folder); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating folder",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": folder,
	})
}

func (h *PhotoFolderHandler) UpdatePhotoFolder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid photo folder ID",
			"error":   err.Error(),
		})
	}

	var folder entity.PhotoFolder
	if err := c.BodyParser(&folder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}

	err = h.usecase.UpdatePhotoFolder(&folder, uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating photo folder",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Photo folder updated successfully",
		"data":    folder,
	})
}

func (h *PhotoFolderHandler) DeletePhotoFolder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid photo folder ID",
			"error":   err.Error(),
		})
	}

	if err := h.usecase.DeletePhotoFolder(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting photo folder",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Photo folder deleted successfully",
	})
}
