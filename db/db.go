package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/islombay/cms/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("db.db"), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Silent),
		TranslateError: true,
	})
	if err != nil {
		panic("could not open db")
	}

	if err := db.AutoMigrate(&model.User{}, &model.Task{}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create gen pod
	var user model.User
	db.First(&user, "username = ?", "root")

	if user.ID == "" {
		db.Create(&model.User{
			ID:       uuid.NewString(),
			Username: "root",
			Password: "root",
			Role:     string(model.UserGenPod),
		})
	}

	// create qa
	var qa model.User
	db.First(&qa, "username = ?", "qa")

	if qa.ID == "" {
		db.Create(&model.User{
			ID:       uuid.NewString(),
			Username: "qa",
			Password: "qa",
			Role:     string(model.UserQA),
		})
	}

	return db
}
