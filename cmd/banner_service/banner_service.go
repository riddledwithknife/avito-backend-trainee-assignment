package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"avito-backend-trainee-assignment/internal/models"
)

func LoadConfig(filename string) models.Config {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := models.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalf("Failed to decode config file: %v", err)
	}

	return config
}

func getUserBannerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
		if err != nil {
			http.Error(w, "Invalid tag ID", http.StatusBadRequest)
			return
		}

		featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
		if err != nil {
			http.Error(w, "Invalid feature ID", http.StatusBadRequest)
			return
		}

		_, err = strconv.ParseBool(r.URL.Query().Get("use_last_revision")) // Use last revision mechanic to add

		token := r.Header.Get("token")
		var compare models.Token
		err = db.Where("token = ?", token).First(&compare).Error //Если админ, то выводить неактивные
		if token == "" || err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var banner models.Banner
		if err := db.Where("feature_id = ? AND ? = ANY(tag_ids)", featureID, tagID).First(&banner).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Banner not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to get banner", http.StatusInternalServerError)
			}
			return
		}

		if !banner.IsActive {
			http.Error(w, "Banner is not active", http.StatusNotFound)
			return
		}

		bannerJSON, err := json.Marshal(banner)
		if err != nil {
			http.Error(w, "Failed to marshal banner to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bannerJSON)
	}
}

func getBannersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var adminToken models.Token
		err := db.Where("token = ? AND is_admin = true", token).First(&adminToken).Error
		if err != nil {
			http.Error(w, "Not an admin", http.StatusForbidden)
			return
		}

		featureIDStr := r.URL.Query().Get("feature_id")
		tagIDStr := r.URL.Query().Get("tag_id")
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		var featureID, tagID, limit, offset int

		if featureIDStr != "" {
			featureID, err = strconv.Atoi(featureIDStr)
			if err != nil {
				http.Error(w, "Invalid feature ID", http.StatusInternalServerError)
				return
			}
		}

		if tagIDStr != "" {
			tagID, err = strconv.Atoi(tagIDStr)
			if err != nil {
				http.Error(w, "Invalid tag ID", http.StatusInternalServerError)
				return
			}
		}

		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				http.Error(w, "Invalid limit", http.StatusInternalServerError)
				return
			}
		}

		if offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				http.Error(w, "Invalid offset", http.StatusInternalServerError)
				return
			}
		}

		var banners []models.Banner
		query := db.Model(&models.Banner{})
		if featureID != 0 {
			query = query.Where("feature_id = ?", featureID)
		}
		if tagID != 0 {
			query = query.Where("? = ANY(tag_ids)", tagID)
		}
		if limit != 0 {
			query = query.Limit(limit)
		}
		if offset != 0 {
			query = query.Offset(offset)
		}
		if err := query.Find(&banners).Error; err != nil {
			http.Error(w, "Failed to fetch banners", http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(banners)
		if err != nil {
			http.Error(w, "Failed to marshal banners to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func createBannerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var adminToken models.Token
		err := db.Where("token = ? AND is_admin = true", token).First(&adminToken).Error
		if err != nil {
			http.Error(w, "Not an admin", http.StatusForbidden)
			return
		}

		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()

		var newBanner models.Banner
		if err := decoder.Decode(&newBanner); err != nil {
			http.Error(w, "Failed to decode request body", http.StatusBadRequest)
			return
		}

		if err := db.Create(&newBanner).Error; err != nil {
			http.Error(w, "Failed to create banner", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":` + strconv.Itoa(int(newBanner.ID)) + `}`))
	}
}

func updateBannerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var adminToken models.Token
		err := db.Where("token = ? AND is_admin = true", token).First(&adminToken).Error
		if err != nil {
			http.Error(w, "Not an admin", http.StatusForbidden)
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid banner ID", http.StatusBadRequest)
			return
		}

		var banner models.Banner
		if err := db.First(&banner, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Banner not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to fetch banner", http.StatusInternalServerError)
			}
			return
		}

		var updatedBanner models.Banner
		if err := json.NewDecoder(r.Body).Decode(&updatedBanner); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		banner.FeatureID = updatedBanner.FeatureID
		banner.Title = updatedBanner.Title
		banner.Text = updatedBanner.Text
		banner.URL = updatedBanner.URL
		banner.IsActive = updatedBanner.IsActive

		if err := db.Save(&banner).Error; err != nil {
			http.Error(w, "Failed to update banner", http.StatusInternalServerError) //Выводить жисон с ошибкой, возможно на всех 500х
			return
		}

		response, err := json.Marshal(banner)
		if err != nil {
			http.Error(w, "Failed to marshal updated banner to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json") //Проверить ответы
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func deleteBannerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		bannerID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid banner ID", http.StatusBadRequest)
			return
		}

		var banner models.Banner
		result := db.First(&banner, bannerID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Banner not found", http.StatusNotFound)
			return
		} else if result.Error != nil {
			http.Error(w, "Failed to delete banner", http.StatusInternalServerError)
			return
		}

		if err := db.Delete(&banner).Error; err != nil {
			http.Error(w, "Failed to delete banner", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func createInitUsersHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := models.Token{
			Token:   "user",
			IsAdmin: false,
		}
		admin := models.Token{
			Token:   "admin",
			IsAdmin: true,
		}
		db.Create(&user)
		db.Create(&admin)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"user":"created", "admin":"created"}`))
	}
}

func main() {
	config := LoadConfig("./config.json")
	log.Println("Config file loaded")

	db, err := gorm.Open(postgres.Open("host="+config.Database.Host+" user="+config.Database.User+" password="+config.Database.Password+
		" dbname="+config.Database.DBName+" port="+config.Database.Port+" sslmode=disable"), &gorm.Config{})
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

	router.HandleFunc("/user_banner", getUserBannerHandler(db)).Methods("GET")
	router.HandleFunc("/banner", getBannersHandler(db)).Methods("GET")
	router.HandleFunc("/banner", createBannerHandler(db)).Methods("POST")
	router.HandleFunc("/banner/{id}", updateBannerHandler(db)).Methods("PATCH")
	router.HandleFunc("/banner/{id}", deleteBannerHandler(db)).Methods("DELETE")
	router.HandleFunc("/createInitUsers", createInitUsersHandler(db))

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
