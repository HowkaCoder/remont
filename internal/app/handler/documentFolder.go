package handler

import (
	"strconv"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
)

type DocumentFolderHandler struct {
	usecase usecase.DocumentFolderUsecase
}

func NewDocumentFolderHandler(usecase usecase.DocumentFolderUsecase) *DocumentFolderHandler {
	return &DocumentFolderHandler{usecase: usecase}
}
func (h *DocumentFolderHandler) GetAllDocumentFolders(c *fiber.Ctx) error {
	folders, err := h.usecase.GetAllDocumentFolders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting document folders",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": folders,
	})
}

func (h *DocumentFolderHandler) GetDocumentFolderByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid document folder ID",
			"error":   err.Error(),
		})
	}

	folder, err := h.usecase.GetDocumentFolderByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting document folder",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": folder,
	})
}

func (h *DocumentFolderHandler) GetDocumentFoldersByProjectID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid project ID",
			"error":   err.Error(),
		})
	}

	folders, err := h.usecase.GetDocumentFoldersByProjectID(uint(id))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error getting document folders for project",
			"error":   err.Error(),
		})
	}

	return c.JSON(folders)
}

func (h *DocumentFolderHandler) CreateDocumentFolder(c *fiber.Ctx) error {
	var folder entity.DocumentFolder
	if err := c.BodyParser(&folder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}

	err := h.usecase.CreateDocumentFolder(&folder)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error creating document folder",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Document folder created successfully",
	})
}

func (h *DocumentFolderHandler) UpdateDocumentFolder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid document folder ID",
			"error":   err.Error(),
		})
	}

	var folder *entity.DocumentFolder
	if err := c.BodyParser(&folder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request payload",
			"error":   err.Error(),
		})
	}

	err = h.usecase.UpdateDocumentFolder(folder, uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error updating document folder",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Document folder updated successfully",
		"data":    folder,
	})
}

func (h *DocumentFolderHandler) DeleteDocumentFolder(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid document folder ID",
			"error":   err.Error(),
		})
	}

	err = h.usecase.DeleteDocumentFolder(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Error deleting document folder",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Document folder deleted successfully",
	})
}
