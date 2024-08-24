package test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/HowkaCoder/remont/internal"
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/handler"
	"github.com/HowkaCoder/remont/internal/app/repository"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func SetupApp() *fiber.App {
	app := fiber.New()

	db := internal.Init()

	docRepo := repository.NewDocumentRepository(db)
	docUsecase := usecase.NewDocumentUsecase(docRepo)
	docHandler := handler.NewDocumentHandler(docUsecase)

	app.Get("/api/docs", docHandler.GetAllDocuments)
	app.Post("/api/docs", docHandler.CreateDocument)
	app.Patch("/api/docs/:id", docHandler.UpdateDocument)
	app.Get("/api/docs/:id", docHandler.GetDocumentByID)
	app.Get("/api/docs/folder/:folder_id", docHandler.GetDocumentsByFolderID)
	app.Delete("/api/docs/:id", docHandler.DeleteDocument)

	return app
}

func TestCreateDocumentMultipart(t *testing.T) {
	app := SetupApp()

	// Создание buffer и writer для multipart формы
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Добавление полей формы
	_ = writer.WriteField("name", "Document 1")
	_ = writer.WriteField("projectID", "1")
	_ = writer.WriteField("documentFolderID", "1")

	// Добавление файла
	file, err := os.Open("test_data/photo_1_2024-08-24_14-30-57.jpg")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", "photo_1_2024-08-24_14-30-57.jpg")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()

	// Создание HTTP-запроса
	req := httptest.NewRequest(http.MethodPost, "/documents", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполнение запроса
	resp, _ := app.Test(req, -1)

	// Проверка ответа
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var doc entity.Document
	json.NewDecoder(resp.Body).Decode(&doc)
	assert.Equal(t, "Document 1", doc.Name)
}
