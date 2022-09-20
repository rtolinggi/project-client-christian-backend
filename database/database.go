package database

import (
	"fmt"
	"log"
	"strconv"

	"github.com/rtolinggi/sales-api/config"
	// "github.com/rtolinggi/sales-api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "gorm.io/gorm/logger"
)

type DBinstance struct {
	DB *gorm.DB
}

var DB DBinstance

func ConnectDB() {
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		log.Fatal("Failed inital Port \n", err)
	}

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Config("DB_HOST"), port, config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		log.Fatal("Failed to connect database \n", err)
	}

	log.Println("Database Connected")
	// db.Logger = logger.Default.LogMode(logger.Info)
	// log.Println("running migrations")
	// db.AutoMigrate(&models.Karyawan{}, &models.User{})

	DB = DBinstance{
		DB: db,
	}

}
