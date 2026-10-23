package main

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/HowkaCoder/remont/internal"
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/config"
	"github.com/gofiber/fiber/v2"
)

const testPassword = "password123"

func setup(t *testing.T) (*fiber.App, config.Config) {
	t.Helper()
	dir := t.TempDir()
	cfg := config.Config{
		DBPath:     filepath.Join(dir, "test.db"),
		UploadsDir: filepath.Join(dir, "uploads"),
		PublicURL:  "http://test",
		JWTSecret:  "test-secret",
	}

	db, err := internal.Open(cfg.DBPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := internal.SeedRoles(db); err != nil {
		t.Fatal(err)
	}
	if err := internal.SeedUser(db, entity.User{Email: "admin@test.local"}, testPassword, "admin"); err != nil {
		t.Fatal(err)
	}
	if err := internal.SeedUser(db, entity.User{Email: "manager@test.local"}, testPassword, "manager"); err != nil {
		t.Fatal(err)
	}

	return newApp(db, cfg), cfg
}

func request(t *testing.T, app *fiber.App, method, url, token string, body io.Reader, contentType string) (int, map[string]any) {
	t.Helper()
	req := httptest.NewRequest(method, url, body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	out := map[string]any{}
	json.NewDecoder(resp.Body).Decode(&out)
	return resp.StatusCode, out
}

func jsonBody(v any) io.Reader {
	b, _ := json.Marshal(v)
	return bytes.NewReader(b)
}

func login(t *testing.T, app *fiber.App, email, password string) (int, string) {
	t.Helper()
	status, body := request(t, app, http.MethodPost, "/login", "", jsonBody(map[string]string{"email": email, "password": password}), "application/json")
	token, _ := body["token"].(string)
	return status, token
}

func TestLogin(t *testing.T) {
	app, _ := setup(t)

	if status, token := login(t, app, "admin@test.local", testPassword); status != http.StatusOK || token == "" {
		t.Fatalf("valid login: status %d", status)
	}
	if status, _ := login(t, app, "admin@test.local", "wrong-password"); status != http.StatusUnauthorized {
		t.Fatalf("wrong password: expected 401, got %d", status)
	}
}

func TestRegisterCreatesClient(t *testing.T) {
	app, _ := setup(t)

	status, _ := request(t, app, http.MethodPost, "/register", "", jsonBody(map[string]string{
		"FirstName": "Ali",
		"Email":     "ali@test.local",
		"Password":  testPassword,
	}), "application/json")
	if status != http.StatusCreated {
		t.Fatalf("register: %d", status)
	}

	_, token := login(t, app, "ali@test.local", testPassword)
	if status, _ := request(t, app, http.MethodGet, "/api/projects", token, nil, ""); status != http.StatusOK {
		t.Fatalf("client should read projects, got %d", status)
	}
	if status, _ := request(t, app, http.MethodPost, "/roles/assign", token, jsonBody(map[string]uint{"RoleID": 1, "PermissionID": 1}), "application/json"); status != http.StatusForbidden {
		t.Fatalf("client must not manage roles, got %d", status)
	}
}

func TestRoutesRequireToken(t *testing.T) {
	app, _ := setup(t)
	for _, path := range []string{"/users", "/roles", "/api/projects", "/get-profile"} {
		if status, _ := request(t, app, http.MethodGet, path, "", nil, ""); status != http.StatusUnauthorized {
			t.Errorf("%s: expected 401, got %d", path, status)
		}
	}
}

func TestPasswordNotExposed(t *testing.T) {
	app, _ := setup(t)
	_, token := login(t, app, "admin@test.local", testPassword)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("admin list users: %d", resp.StatusCode)
	}
	if bytes.Contains(raw, []byte("Password")) || bytes.Contains(raw, []byte("$2a$")) {
		t.Fatal("password hash is exposed")
	}
}

func TestUploadDocument(t *testing.T) {
	app, cfg := setup(t)
	_, token := login(t, app, "manager@test.local", testPassword)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("name", "contract")
	w.WriteField("projectID", "1")
	w.WriteField("folderID", "1")
	fw, _ := w.CreateFormFile("file", "../../evil.pdf")
	fw.Write([]byte("%PDF-1.4"))
	w.Close()

	status, body := request(t, app, http.MethodPost, "/api/docs", token, &buf, w.FormDataContentType())
	if status != http.StatusCreated {
		t.Fatalf("upload: status %d, body %v", status, body)
	}

	entries, err := os.ReadDir(filepath.Join(cfg.UploadsDir, "documents"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 || filepath.Ext(entries[0].Name()) != ".pdf" {
		t.Fatalf("file is not stored inside uploads/documents: %v", entries)
	}

	id := uint(body["id"].(float64))
	status, _ = request(t, app, http.MethodDelete, "/api/docs/"+itoa(id), token, nil, "")
	if status != http.StatusOK {
		t.Fatalf("delete: %d", status)
	}
	entries, _ = os.ReadDir(filepath.Join(cfg.UploadsDir, "documents"))
	if len(entries) != 0 {
		t.Fatal("file was not removed from disk")
	}
}

func TestUploadRequiresFields(t *testing.T) {
	app, _ := setup(t)
	_, token := login(t, app, "manager@test.local", testPassword)

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("name", "contract")
	w.Close()

	if status, _ := request(t, app, http.MethodPost, "/api/docs", token, &buf, w.FormDataContentType()); status != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing fields, got %d", status)
	}
}

func itoa(n uint) string {
	b, _ := json.Marshal(n)
	return string(b)
}
