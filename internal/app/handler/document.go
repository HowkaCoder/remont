package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
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

		eDoc, err := dh.usecadse.GetDocumentByID(uint(id))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
		}

		if err := os.Remove(eDoc.Filepath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
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

		doc := entity.Document{ProjectID: uint(projectID), Name: name[0], Filepath: filePath}

		if err := dh.usecadse.UpdateDocument(&doc, uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
		}

	} else {

		doc := entity.Document{ProjectID: uint(projectID), Name: name[0]}

		if err := dh.usecadse.UpdateDocument(&doc, uint(id)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
		}

	}
	return c.JSON(fiber.Map{"name": name[0], "projectID": projectID, "file": files, "len(files)": len(files), "id": uint(id)})
}

func (dh *DocumentHandler) DeleteDocument(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}

	doc, err := dh.usecadse.GetDocumentByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
	}

	if err := os.Remove(doc.Filepath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
	}

	if err := dh.usecadse.DeleteDocument(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "successfully deleted"})

}

func (dh *DocumentHandler) GetDocumentsByFolderID(c *fiber.Ctx) error {
	folderID, err := strconv.Atoi(c.Params("folderID"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"Error": err.Error()})
	}
	docs, err := dh.usecadse.GetDocumentsByFolderID(uint(folderID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"Error": err.Error()})
	}

	return c.JSON(docs)
}
