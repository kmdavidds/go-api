package initializers

import "github.com/kmdavidds/go-api/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Scavenger{})
	DB.AutoMigrate(&models.Agenda{})
	DB.AutoMigrate(&models.TakenAgenda{})
	DB.AutoMigrate(&models.Petan{})
}
