package main

import (
	"log"

	"github.com/HowkaCoder/remont/internal"
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/handler"
	"github.com/HowkaCoder/remont/internal/app/repository"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var db *gorm.DB
var jwtSecret = []byte("your_secret_key")

func main() {
	// Инициализация базы данных
	db = internal.Init()

	// Инициализация репозиториев, usecase и обработчиков
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo)
	userHandler := handler.NewUserHandler(userUsecase)

	docRepo := repository.NewDocumentRepository(db)
	docUsecase := usecase.NewDocumentUsecase(docRepo)
	docHandler := handler.NewDocumentHandler(docUsecase)

	projectRepo := repository.NewProjectRepository(db)
	projectUsecase := usecase.NewProjectUsecase(projectRepo)
	projectHandler := handler.NewProjectHandler(projectUsecase)

	app := fiber.New()

	// Роуты для регистрации и логина
	app.Post("/register", userHandler.CreateUser)
	app.Post("/login", login)

	app.Post("/roles", createRole)
	app.Post("/permissions", createPermission)
	app.Post("/roles/assign", assignPermissionToRole)
	app.Post("/users/assign-role", assignRoleToUser)

	// Middleware для аутентификации
	app.Use(authMiddleware)

	// Роуты для управления документами с проверкой прав доступа
	app.Get("/api/docs", PermissionMiddleware("get_documents"), docHandler.GetAllDocuments)
	app.Post("/api/docs", PermissionMiddleware("create_document"), docHandler.CreateDocument)
	app.Patch("/api/docs/:id", PermissionMiddleware("update_document"), docHandler.UpdateDocument)
	app.Get("/api/docs/:id", PermissionMiddleware("get_document_by_id"), docHandler.GetDocumentByID)
	app.Delete("/api/docs/:id", PermissionMiddleware("delete_document"), docHandler.DeleteDocument)

	app.Get("/api/projects", PermissionMiddleware("get_all_projects"), projectHandler.GetAllProjects)
	app.Post("/api/projects", PermissionMiddleware("create_project"), projectHandler.CreateProject)
	app.Patch("/api/projects/:id", PermissionMiddleware("update_project"), projectHandler.UpdateProject)
	app.Get("/api/projects/:id", PermissionMiddleware("get_project_by_id"), projectHandler.GetProjectByID)
	app.Delete("/api/projects/:id", PermissionMiddleware("delete_project"), projectHandler.DeleteProject)

	app.Get("/api/project-roles", PermissionMiddleware("get_all_project_roles"), projectHandler.GetAllProjectRole)
	app.Post("/api/project-roles", PermissionMiddleware("create_project_role"), projectHandler.CreateProjectRole)
	app.Patch("/api/project-roles/:id", PermissionMiddleware("update_project_role"), projectHandler.UpdateProjectRole)
	app.Get("/api/project-roles/:id", PermissionMiddleware("get_project_role_by_id"), projectHandler.GetProjectRoleByID)
	app.Delete("/api/project-roles/:id", PermissionMiddleware("delete_project_role"), projectHandler.DeleteProjectRole)

	// Роуты для управления ролями и правами
	/*app.Post("/roles", PermissionMiddleware("create_role"), createRole)
	app.Post("/permissions", PermissionMiddleware("create_permission"), createPermission)
	app.Post("/roles/assign", PermissionMiddleware("assign_permission"), assignPermissionToRole)
	app.Post("/users/assign-role", PermissionMiddleware("assign_role"), assignRoleToUser)
	*/
	log.Fatal(app.Listen(":3000"))
}

// Middleware для проверки прав доступа
func PermissionMiddleware(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("userID").(uint) // Предполагается, что userID установлен в контексте после аутентификации
		var user entity.User
		db.Preload("Roles.Permissions").First(&user, userID)

		for _, role := range user.Roles {
			for _, permission := range role.Permissions {
				if permission.Name == requiredPermission {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden",
		})
	}
}

// Middleware для аутентификации
func authMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or malformed JWT"})
	}

	tokenString := authHeader[len("Bearer "):]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired JWT"})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["userID"].(float64))
	c.Locals("userID", userID)

	return c.Next()
}

// Функция для логина
func login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user entity.User
	db.Where("email = ?", req.Email).First(&user)
	if user.ID == 0 || user.Password != req.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"token": tokenString})
}

// Функция для создания роли
func createRole(c *fiber.Ctx) error {
	var role entity.Role
	if err := c.BodyParser(&role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	db.Create(&role)
	return c.Status(fiber.StatusCreated).JSON(role)
}

// Функция для создания разрешения
func createPermission(c *fiber.Ctx) error {
	var permission entity.Permission
	if err := c.BodyParser(&permission); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	db.Create(&permission)
	return c.Status(fiber.StatusCreated).JSON(permission)
}

// Функция для назначения разрешения роли
func assignPermissionToRole(c *fiber.Ctx) error {
	var input struct {
		RoleID       uint `json:"RoleID"`
		PermissionID uint `json:"PermissionID"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	rolePermission := entity.RolePermission{
		RoleID:       input.RoleID,
		PermissionID: input.PermissionID,
	}

	db.Create(&rolePermission)
	return c.Status(fiber.StatusCreated).JSON(rolePermission)
}

// Функция для назначения роли пользователю
func assignRoleToUser(c *fiber.Ctx) error {
	var input struct {
		UserID uint `json:"UserID"`
		RoleID uint `json:"RoleID"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userRole := entity.UserRole{
		UserID: input.UserID,
		RoleID: input.RoleID,
	}

	db.Create(&userRole)
	return c.Status(fiber.StatusCreated).JSON(userRole)
}
