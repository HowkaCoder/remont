package handler

import (
	"fmt"
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type DocumentHandler struct {
	usecadse usecase.DocumentUsecase
}

func NewDocumentHandler(usecase usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{usecadse: usecase}
}

func (dh *DocumentHandler) GetAllDocuments(c *fiber.Ctx) error {
	docs, err := dh.usecadse.GetAllDocuments()
	if err != nil {
		return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"Error": err.Error()})
	}

	return c.JSON(docs)
}

func (dh *DocumentHandler) GetDocumentByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}
	doc, err := dh.usecadse.GetDocumentByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
	}

	return c.JSON(doc)
}

func (dh *DocumentHandler) CreateDocument(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	name := form.Value["name"]
	if len(name) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	}

	projectID, err := strconv.Atoi(form.Value["projectID"][0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Project ID is required"})
	}

	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File is required"})
	}

	file := files[0]

	uploadDir := "./uploads/documents"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	doc := entity.Document{ProjectID: uint(projectID), Name: name[0], Filepath: uploadDir + "/" + fileName}

	if err := dh.usecadse.CreateDocument(&doc); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
	}

	return c.JSON(doc)
}

func (dh *DocumentHandler) UpdateDocument(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	name := form.Value["name"]

	projectID, err := strconv.Atoi(form.Value["projectID"][0])
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Project ID is required"})
	}

	files := form.File["file"]
	if len(files) > 0 {

		file := files[0]

		uploadDir := "./uploads/documents"
		if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		filePath := filepath.Join(uploadDir, fileName)

		if err := c.SaveFile(file, filePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		doc := entity.Document{ProjectID: uint(projectID), Name: name[0], Filepath: uploadDir + "/" + fileName}

		if err := dh.usecadse.UpdateDocument(&doc, uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
		}

	}

	doc := entity.Document{ProjectID: uint(projectID), Name: name[0]}

	if err := dh.usecadse.UpdateDocument(&doc, uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
	}
	return c.JSON(doc)
}
