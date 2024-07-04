package internal

import (
	"log"

	"github.com/HowkaCoder/remont/internal/app/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() *gorm.DB {
	var err error
	//dsn := "root:pBaYLMKHIVQFHPBbbRKAfphLmzReYKSu@tcp(roundhouse.proxy.rlwy.net:39674)/railway?charset=utf8mb4&parseTime=True&loc=Local"
	//DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	DB, err := gorm.Open(sqlite.Open("database/database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = DB.AutoMigrate(&entity.User{} , &entity.Role{} , &entity.UserRole{} , &entity.Permission{} , &entity.RolePermission{} , &entity.UserPermission{})
	if err != nil {
		log.Fatal(err)
	}
	return DB

	
	return DB
}
