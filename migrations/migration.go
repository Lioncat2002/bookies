package main

import (
	"backend/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}
	dburi := os.Getenv("DB_URI") //used a cockroachdb database but postgres is fine
	db, err := gorm.Open(postgres.Open(dburi), &gorm.Config{})
	if err != nil {
		log.Println("Error coonection to db", err)
	}
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Book{})
	db.Exec("DROP TABLE user_owns")
	db.Exec("DROP TABLE user_carts")
	db.AutoMigrate(&models.User{}, &models.UserCart{}, &models.UserOwns{})
	db.AutoMigrate(&models.Book{})
	log.Println("Successfully migrated")
}
