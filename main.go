package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/HowkaCoder/remont/internal"
	"github.com/HowkaCoder/remont/internal/app/entity"
	"github.com/HowkaCoder/remont/internal/app/handler"
	"github.com/HowkaCoder/remont/internal/app/repository"
	"github.com/HowkaCoder/remont/internal/app/usecase"
	"github.com/HowkaCoder/remont/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := internal.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := seed(db, cfg); err != nil {
		log.Fatalf("seed: %v", err)
	}

	app := newApp(db, cfg)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		if err := app.ShutdownWithTimeout(10 * time.Second); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}()

	log.Printf("server is running on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}

func seed(db *gorm.DB, cfg config.Config) error {
	if err := internal.SeedRoles(db); err != nil {
		return err
	}
	if cfg.AdminEmail != "" && cfg.AdminPassword != "" {
		admin := entity.User{FirstName: "Admin", Email: cfg.AdminEmail}
		if err := internal.SeedUser(db, admin, cfg.AdminPassword, "admin"); err != nil {
			return err
		}
	}
	if cfg.SeedDemo {
		return internal.SeedDemoUsers(db, cfg.DemoPassword)
	}
	return nil
}

func newApp(db *gorm.DB, cfg config.Config) *fiber.App {
	photosDir := filepath.Join(cfg.UploadsDir, "photos")
	documentsDir := filepath.Join(cfg.UploadsDir, "documents")

	auth := handler.NewAuthHandler(repository.NewAuthRepository(db), []byte(cfg.JWTSecret))
	userHandler := handler.NewUserHandler(usecase.NewUserUsecase(repository.NewUserRepository(db)))
	docHandler := handler.NewDocumentHandler(usecase.NewDocumentUsecase(repository.NewDocumentRepository(db)), documentsDir)
	projectHandler := handler.NewProjectHandler(usecase.NewProjectUsecase(repository.NewProjectRepository(db)))
	photoHandler := handler.NewPhotoHandler(usecase.NewPhotoUsecase(repository.NewPhotoRepository(db)), photosDir, cfg.PublicURL)
	documentFolderHandler := handler.NewDocumentFolderHandler(usecase.NewDocumentFolderUsecase(repository.NewDocumentFolderRepository(db)))
	charHandler := handler.NewCharHandler(usecase.NewCharUsecase(repository.NewCharRepository(db)))
	photoFolderHandler := handler.NewPhotoFolderHandler(usecase.NewPhotoFolderUsecase(repository.NewPhotoFolderRepository(db)))
	stateHandler := handler.NewStateHandler(usecase.NewStateUsecase(repository.NewStateRepository(db)))

	app := fiber.New(fiber.Config{BodyLimit: 100 * 1024 * 1024})
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,PATCH,DELETE",
		AllowHeaders: "Content-Type, Authorization",
	}))

	app.Static("/uploads/photos", photosDir)
	app.Static("/uploads/documents", documentsDir)

	app.Post("/register", auth.Register)
	app.Post("/login", auth.Login)

	api := app.Group("", auth.Authenticate)
	can := auth.Require

	api.Get("/get-profile", auth.Profile)

	api.Get("/users", can(internal.PermManageUsers), userHandler.GetAllUsers)
	api.Get("/users/:id", can(internal.PermManageUsers), userHandler.GetUserByID)
	api.Delete("/users/:id", can(internal.PermManageUsers), userHandler.DeleteUser)
	api.Post("/users/assign-role", can(internal.PermManageRoles), auth.AssignRoleToUser)

	api.Get("/roles", can(internal.PermManageRoles), auth.GetRoles)
	api.Post("/roles", can(internal.PermManageRoles), auth.CreateRole)
	api.Post("/roles/assign", can(internal.PermManageRoles), auth.AssignPermissionToRole)
	api.Get("/permissions", can(internal.PermManageRoles), auth.GetPermissions)
	api.Post("/permissions", can(internal.PermManageRoles), auth.CreatePermission)

	api.Get("/api/projects/users/:id", can("get_project_by_id"), auth.ProjectMembers)

	api.Get("/api/docs", can("get_documents"), docHandler.GetAllDocuments)
	api.Post("/api/docs", can("create_document"), docHandler.CreateDocument)
	api.Patch("/api/docs/:id", can("update_document"), docHandler.UpdateDocument)
	api.Get("/api/docs/:id", can("get_document_by_id"), docHandler.GetDocumentByID)
	api.Get("/api/docs/folder/:folderID", can("get_documents_by_folder_id"), docHandler.GetDocumentsByFolderID)
	api.Delete("/api/docs/:id", can("delete_document"), docHandler.DeleteDocument)

	api.Get("/api/projects", can("get_all_projects"), projectHandler.GetAllProjects)
	api.Post("/api/projects", can("create_project"), projectHandler.CreateProject)
	api.Patch("/api/projects/:id", can("update_project"), projectHandler.UpdateProject)
	api.Get("/api/projects/:id", can("get_project_by_id"), projectHandler.GetProjectByID)
	api.Delete("/api/projects/:id", can("delete_project"), projectHandler.DeleteProject)
	api.Get("/api/project/workers/:id", can("get_projects_by_worker_id"), projectHandler.GetAllProjectsAsAWorker)
	api.Get("/api/project/clients/:id", can("get_projects_by_client_id"), projectHandler.GetAllProjectsAsAClient)

	api.Get("/api/project-roles", can("get_all_project_roles"), projectHandler.GetAllProjectRole)
	api.Post("/api/project-roles", can("create_project_role"), projectHandler.CreateProjectRole)
	api.Patch("/api/project-roles/:id", can("update_project_role"), projectHandler.UpdateProjectRole)
	api.Get("/api/project-roles/:id", can("get_project_role_by_id"), projectHandler.GetProjectRoleByID)
	api.Delete("/api/project-roles/:id", can("delete_project_role"), projectHandler.DeleteProjectRole)

	api.Get("/api/photos", can("get_all_photos"), photoHandler.GetAllPhotos)
	api.Post("/api/photos", can("create_photo"), photoHandler.CreatePhoto)
	api.Get("/api/photos/folder/:folderID", can("get_photos_by_folder_id"), photoHandler.GetPhotosByFolderID)
	api.Get("/api/photos/:id", can("get_photo_by_id"), photoHandler.GetPhotoByID)
	api.Patch("/api/photos/:id", can("update_photo"), photoHandler.UpdatePhoto)
	api.Delete("/api/photos/:id", can("delete_photo"), photoHandler.DeletePhoto)

	api.Get("/api/document-folders", can("get_all_document_folders"), documentFolderHandler.GetAllDocumentFolders)
	api.Get("/api/document-folders/:id", can("get_document_folder_by_folder_id"), documentFolderHandler.GetDocumentFolderByID)
	api.Get("/api/document-folders/project/:id", can("get_document_folder_by_project_id"), documentFolderHandler.GetDocumentFoldersByProjectID)
	api.Post("/api/document-folders", can("create_document_folder"), documentFolderHandler.CreateDocumentFolder)
	api.Patch("/api/document-folders/:id", can("update_document_folder"), documentFolderHandler.UpdateDocumentFolder)
	api.Delete("/api/document-folders/:id", can("delete_document_folder"), documentFolderHandler.DeleteDocumentFolder)

	api.Get("/api/photo-folder", can("get_all_photo_folders"), photoFolderHandler.GetAllPhotoFolders)
	api.Get("/api/photo-folder/:id", can("get_photo_folder_by_folder_id"), photoFolderHandler.GetPhotoFolderByID)
	api.Post("/api/photo-folder", can("create_photo_folder"), photoFolderHandler.CreatePhotoFolder)
	api.Patch("/api/photo-folder/:id", can("update_photo_folder"), photoFolderHandler.UpdatePhotoFolder)
	api.Delete("/api/photo-folder/:id", can("delete_photo_folder"), photoFolderHandler.DeletePhotoFolder)
	api.Get("/api/photo-folder/project/:id", can("get_photo_folder_by_project_id"), photoFolderHandler.GetPhotoFoldersByProjectID)

	api.Get("/api/chars/project/:id", can("get_chars_by_project_id"), charHandler.GetAllCharsByProjectID)
	api.Post("/api/chars/project", can("create_char"), charHandler.CreateChar)
	api.Patch("/api/chars/:id", can("update_char"), charHandler.UpdateChar)
	api.Delete("/api/chars/:id", can("delete_char"), charHandler.DeleteChar)

	api.Get("/api/states/:id", can("get_states_by_project_id"), stateHandler.GetStatesByProjectID)
	api.Get("/api/states/worker/:id", can("get_states_by_worker_id"), stateHandler.GetStatesByWorkerID)
	api.Post("/api/states", can("create_state"), stateHandler.CreateState)
	api.Patch("/api/states/:id", can("update_state"), stateHandler.UpdateState)
	api.Delete("/api/states/:id", can("delete_state"), stateHandler.DeleteState)
	api.Post("/api/states/add-user", can("create_state_relation"), stateHandler.AssignWorkerToState)
	api.Post("/api/states/remove-user", can("delete_state_relation"), stateHandler.RemoveWorkerFromState)

	api.Post("/api/details", can("create_detail"), stateHandler.CreateRepairDetails)
	api.Get("/api/details/project/:id", can("get_detail_by_id"), stateHandler.GetRepairDetailsByProjectID)
	api.Patch("/api/details", can("update_detail"), stateHandler.UpdateRepairDetail)
	api.Delete("/api/details/:id", can("delete_detail"), stateHandler.DeleteRepairDetail)

	return app
}
