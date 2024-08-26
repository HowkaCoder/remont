package test

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/HowkaCoder/remont/internal"
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/handler"
	"github.com/HowkaCoder/remont/internal/app/repository"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func TestGetDocument(t *testing.T) {
	app := SetupApp()
	t.Log("TestGetDocument is starting")
	a := time.Now()
	// Creating a document first
	//  doc := entity.Document{Name: "Document 2", Filepath: "/path/to/doc2.pdf", ProjectID: 1, DocumentFolderID: 1}
	//    repository.CreateDocument(&doc)
	name := uuid.New().String()
	// Создание buffer и writer для multipart формы
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Добавление полей формы
	_ = writer.WriteField("name", name)
	_ = writer.WriteField("projectID", "1")
	_ = writer.WriteField("documentFolderID", "1")
	_ = writer.WriteField("folderID", "1")

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
	req := httptest.NewRequest(http.MethodPost, "/api/docs", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполнение запроса
	resp, _ := app.Test(req, -1)

	req = httptest.NewRequest(http.MethodGet, "/api/docs/1", nil)
	resp, _ = app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var retrievedDoc entity.Document
	json.NewDecoder(resp.Body).Decode(&retrievedDoc)
	assert.Equal(t, name, retrievedDoc.Name)
	t.Log(time.Since(a))
}

func TestCreateDocumentMultipart(t *testing.T) {
	t.Log("TestCreateDocumentMultipart is starting")
	a := time.Now()
	app := SetupApp()

	// Создание buffer и writer для multipart формы
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	// Добавление полей формы
	_ = writer.WriteField("name", "Document 1")
	_ = writer.WriteField("projectID", "1")
	_ = writer.WriteField("documentFolderID", "1")
	_ = writer.WriteField("folderID", "1")

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
	req := httptest.NewRequest(http.MethodPost, "/api/docs", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполнение запроса
	resp, _ := app.Test(req, -1)

	// Проверка ответа
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var doc entity.Document
	json.NewDecoder(resp.Body).Decode(&doc)
	assert.Equal(t, "Document 1", doc.Name)
	t.Log(time.Since(a))
}

func TestUpdateDocumentFolder(t *testing.T) {
	app := SetupApp()
	a := time.Now()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	name := uuid.New().String()
	// Добавление полей формы
	_ = writer.WriteField("name", name)
	_ = writer.WriteField("projectID", "1")
	_ = writer.WriteField("documentFolderID", "1")
	_ = writer.WriteField("folderID", "1")
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
	req := httptest.NewRequest(http.MethodPost, "/api/docs", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Выполнение запроса
	resp, _ := app.Test(req, -1)

	payload := `{
        "name": "Updated Document 3",
	"projectID": 1,
        "documentFolderID": 1
    }`

	req = httptest.NewRequest(http.MethodPatch, "/api/docs/3", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, _ = app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var updatedDoc entity.Document
	json.NewDecoder(resp.Body).Decode(&updatedDoc)
	t.Log(updatedDoc)
	assert.Equal(t, "Updated Document 3", updatedDoc.Name)

	t.Log(time.Since(a))

}
