package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/norfabagas/auth/api/models"
)

var users = []models.User{
	models.User{
		Name:     "admin",
		Email:    "admin@mail.com",
		Password: "password",
	},
}

func Load(db *gorm.DB) {
	err := db.Debug().DropTableIfExists(&models.User{}).Error
	if err != nil {
		log.Fatalf("Cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}).Error
	if err != nil {
		log.Fatalf("Cannot migrate table: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("Cannot seed users table: %v", err)
		}
	}
}
