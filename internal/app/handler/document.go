package handler

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
)

type DocumentHandler struct {
	usecase usecase.DocumentUsecase
	dir     string
}

func NewDocumentHandler(usecase usecase.DocumentUsecase, dir string) *DocumentHandler {
	return &DocumentHandler{usecase: usecase, dir: dir}
}

func (dh *DocumentHandler) GetAllDocuments(c *fiber.Ctx) error {
	docs, err := dh.usecase.GetAllDocuments()
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(docs)
}

func (dh *DocumentHandler) GetDocumentByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return badRequest(c, errors.New("invalid document id"))
	}

	doc, err := dh.usecase.GetDocumentByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "document not found"})
	}
	return c.JSON(doc)
}

func (dh *DocumentHandler) CreateDocument(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return badRequest(c, err)
	}

	name := formValue(form, "name")
	if name == "" {
		return badRequest(c, errors.New("name is required"))
	}
	projectID, err := formUint(form, "projectID")
	if err != nil {
		return badRequest(c, err)
	}
	folderID, err := formUint(form, "folderID")
	if err != nil {
		return badRequest(c, err)
	}

	file := formFile(form)
	if file == nil {
		return badRequest(c, errors.New("file is required"))
	}

	fileName, err := saveUpload(c, file, dh.dir)
	if err != nil {
		return serverError(c, err)
	}

	doc := entity.Document{
		Name:             name,
		Filepath:         "./uploads/documents/" + fileName,
		ProjectID:        projectID,
		DocumentFolderID: folderID,
		FileExt:          strings.ToLower(filepath.Ext(file.Filename)),
	}
	if err := dh.usecase.CreateDocument(&doc); err != nil {
		removeUpload(dh.dir, fileName)
		return serverError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(doc)
}

func (dh *DocumentHandler) UpdateDocument(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return badRequest(c, errors.New("invalid document id"))
	}

	form, err := c.MultipartForm()
	if err != nil {
		return badRequest(c, err)
	}

	existing, err := dh.usecase.GetDocumentByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "document not found"})
	}

	doc := entity.Document{Name: formValue(form, "name")}
	if v, err := formUint(form, "projectID"); err == nil {
		doc.ProjectID = v
	}
	if v, err := formUint(form, "folderID"); err == nil {
		doc.DocumentFolderID = v
	}

	if file := formFile(form); file != nil {
		fileName, err := saveUpload(c, file, dh.dir)
		if err != nil {
			return serverError(c, err)
		}
		doc.Filepath = "./uploads/documents/" + fileName
		doc.FileExt = strings.ToLower(filepath.Ext(file.Filename))
	}

	if err := dh.usecase.UpdateDocument(&doc, uint(id)); err != nil {
		return serverError(c, err)
	}

	if doc.Filepath != "" {
		removeUpload(dh.dir, existing.Filepath)
	}

	return c.JSON(fiber.Map{"message": "successfully updated"})
}

func (dh *DocumentHandler) DeleteDocument(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return badRequest(c, errors.New("invalid document id"))
	}

	doc, err := dh.usecase.GetDocumentByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "document not found"})
	}

	if err := dh.usecase.DeleteDocument(uint(id)); err != nil {
		return serverError(c, err)
	}
	if err := removeUpload(dh.dir, doc.Filepath); err != nil {
		return serverError(c, err)
	}

	return c.JSON(fiber.Map{"message": "successfully deleted"})
}

func (dh *DocumentHandler) GetDocumentsByFolderID(c *fiber.Ctx) error {
	folderID, err := strconv.Atoi(c.Params("folderID"))
	if err != nil {
		return badRequest(c, errors.New("invalid folder id"))
	}

	docs, err := dh.usecase.GetDocumentsByFolderID(uint(folderID))
	if err != nil {
		return serverError(c, err)
	}
	return c.JSON(docs)
}
