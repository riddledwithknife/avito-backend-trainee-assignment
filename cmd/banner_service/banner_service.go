package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"avito-backend-trainee-assignment/internal/config"
	"avito-backend-trainee-assignment/internal/handlers"
	"avito-backend-trainee-assignment/internal/models"
)

func main() {
	configFile := config.LoadConfig("./config.json")
	log.Println("Config file loaded")

	db, err := gorm.Open(postgres.Open("host="+configFile.Database.Host+" user="+configFile.Database.User+" password="+configFile.Database.Password+
		" dbname="+configFile.Database.DBName+" port="+configFile.Database.Port+" sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect to db: ", err)
	}
	log.Println("Connected to db")

	err = db.AutoMigrate(&models.Banner{}, &models.Token{})
	if err != nil {
		log.Fatalln("Failed to migrate db: ", err)
	}
	log.Println("Database migrated")

	router := mux.NewRouter()

	router.HandleFunc("/user_banner", handlers.GetUserBannerHandler(db)).Methods("GET")
	router.HandleFunc("/banner", handlers.GetBannersHandler(db)).Methods("GET")
	router.HandleFunc("/banner", handlers.CreateBannerHandler(db)).Methods("POST")
	router.HandleFunc("/banner/{id}", handlers.UpdateBannerHandler(db)).Methods("PATCH")
	router.HandleFunc("/banner/{id}", handlers.DeleteBannerHandler(db)).Methods("DELETE")
	router.HandleFunc("/createInitUsers", handlers.CreateInitUsersHandler(db))

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
