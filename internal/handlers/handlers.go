package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"avito-backend-trainee-assignment/internal/models"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func GetUserBannerHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		var isAdmin bool
		if token != "" {
			var compare models.Token
			err := db.Where("token = ?", token).First(&compare).Error
			if err == nil {
				isAdmin = compare.IsAdmin
			}
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
		if err != nil {
			errResponse := ErrorResponse{Error: "Invalid tag id"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
		if err != nil {
			errResponse := ErrorResponse{Error: "Invalid feature id"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		_, err = strconv.ParseBool(r.URL.Query().Get("use_last_revision")) //Not forget

		var banner models.Banner
		if isAdmin {
			err = db.Where("? = ANY(tag_ids) AND feature_id = ?", tagID, featureID).First(&banner).Error
		} else {
			err = db.Where("? = ANY(tag_ids) AND feature_id = ? AND is_active = ?", tagID, featureID, true).First(&banner).Error
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Banner not found", http.StatusNotFound)
			return
		}

		bannerJSON, err := json.Marshal(banner)
		if err != nil {
			errResponse := ErrorResponse{Error: "Error serializing banner"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bannerJSON)
	}
}

func GetBannersHandler(db *gorm.DB) http.HandlerFunc {
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
				errResponse := ErrorResponse{Error: "Invalid feature id"}
				jsonErrResponse, _ := json.Marshal(errResponse)

				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonErrResponse)
				return
			}
		}

		if tagIDStr != "" {
			tagID, err = strconv.Atoi(tagIDStr)
			if err != nil {
				errResponse := ErrorResponse{Error: "Invalid tag id"}
				jsonErrResponse, _ := json.Marshal(errResponse)

				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonErrResponse)
				return
			}
		}

		if limitStr != "" {
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				errResponse := ErrorResponse{Error: "Invalid limit"}
				jsonErrResponse, _ := json.Marshal(errResponse)

				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonErrResponse)
				return
			}
		}

		if offsetStr != "" {
			offset, err = strconv.Atoi(offsetStr)
			if err != nil {
				errResponse := ErrorResponse{Error: "Invalid offset"}
				jsonErrResponse, _ := json.Marshal(errResponse)

				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonErrResponse)
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
		if err = query.Find(&banners).Error; err != nil {
			errResponse := ErrorResponse{Error: "Error fetching banners"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		response, err := json.Marshal(banners)
		if err != nil {
			errResponse := ErrorResponse{Error: "Error serializing banners"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func CreateBannerHandler(db *gorm.DB) http.HandlerFunc {
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
		if err = decoder.Decode(&newBanner); err != nil {
			errResponse := ErrorResponse{Error: "Error deserializing banner"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		if err = db.Create(&newBanner).Error; err != nil { //Не создает false active?
			errResponse := ErrorResponse{Error: "Error creating banner"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":` + strconv.Itoa(int(newBanner.ID)) + `}`))
	}
}

func UpdateBannerHandler(db *gorm.DB) http.HandlerFunc {
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
			errResponse := ErrorResponse{Error: "Invalid id"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		var banner models.Banner
		if err = db.First(&banner, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "Banner not found", http.StatusNotFound)
			} else {
				errResponse := ErrorResponse{Error: "Error fetching banner"}
				jsonErrResponse, _ := json.Marshal(errResponse)

				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonErrResponse)
			}
			return
		}

		var updatedBanner models.Banner
		if err = json.NewDecoder(r.Body).Decode(&updatedBanner); err != nil {
			errResponse := ErrorResponse{Error: "Error deserializing banner"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		banner.FeatureID = updatedBanner.FeatureID
		banner.TagIDs = updatedBanner.TagIDs
		banner.Title = updatedBanner.Title
		banner.Text = updatedBanner.Text
		banner.URL = updatedBanner.URL
		banner.IsActive = updatedBanner.IsActive
		banner.UpdatedAt = time.Now()

		if err = db.Save(&banner).Error; err != nil {
			errResponse := ErrorResponse{Error: "Error updating banner"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteBannerHandler(db *gorm.DB) http.HandlerFunc {
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
		bannerID, err := strconv.Atoi(vars["id"])
		if err != nil {
			errResponse := ErrorResponse{Error: "Invalid id"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		var banner models.Banner
		result := db.First(&banner, bannerID)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Banner not found", http.StatusNotFound)
			return
		} else if result.Error != nil {
			errResponse := ErrorResponse{Error: "Error fetching banner"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		if err = db.Delete(&banner).Error; err != nil {
			errResponse := ErrorResponse{Error: "Error deleting banner"}
			jsonErrResponse, _ := json.Marshal(errResponse)

			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonErrResponse)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func CreateInitUsersHandler(db *gorm.DB) http.HandlerFunc {
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
