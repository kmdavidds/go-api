package initializers

import "github.com/kmdavidds/go-api/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
}
