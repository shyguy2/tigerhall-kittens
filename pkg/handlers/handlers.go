package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"tigerhall-kittens-app/pkg/auth"
	"tigerhall-kittens-app/pkg/models"
	"tigerhall-kittens-app/pkg/service"
	"tigerhall-kittens-app/pkg/utils"
)

const DefaultPageSize = 10

type pagination map[string]interface{}

type handlers struct {
	Auth         *auth.Auth
	Logger       *log.Logger
	TigerService service.TigerService
}

func NewHandlers(tigerService service.TigerService, logger *log.Logger, auth *auth.Auth) *handlers {
	return &handlers{
		Auth:         auth,
		Logger:       logger,
		TigerService: tigerService,
	}
}

func (h *handlers) SignupHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get user data
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	// Validate user data
	if err := auth.ValidateUserData(user); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := h.TigerService.SignupService(&user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Respond with success status
	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"message": "success"})
}

func (h *handlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get login credentials
	var loginCredentials models.LoginCredentials
	if err := json.NewDecoder(r.Body).Decode(&loginCredentials); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	user, err := h.TigerService.LoginService(loginCredentials)
	if err != nil {
		utils.RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	// Generate JWT token
	token, err := h.Auth.GenerateToken(user.Username, user.Email)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Respond with the token as JSON
	response := map[string]string{"token": token}
	utils.RespondWithJSON(w, http.StatusOK, response)
}

func (h *handlers) CreateTigerHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get tiger data
	var tiger models.Tiger
	if err := json.NewDecoder(r.Body).Decode(&tiger); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to parse request body")
		return
	}

	err := h.TigerService.CreateTigerService(tiger)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Respond with success status
	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"message": "success"})
}

func (h *handlers) GetAllTigersHandler(w http.ResponseWriter, r *http.Request) {
	// Get the pagination parameters from the query string
	pageStr := r.FormValue("page")
	pageSizeStr := r.FormValue("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = DefaultPageSize
	}

	tigers, totalCount, err := h.TigerService.GetAllTigersService(page, pageSize)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Construct pagination response
	paginationResponse := pagination{
		"page":       page,
		"pageSize":   pageSize,
		"totalCount": totalCount,
		"totalPages": int(math.Ceil(float64(totalCount) / float64(pageSize))),
		"tigerList":  tigers,
	}

	// Respond with the tigers as JSON
	utils.RespondWithJSON(w, http.StatusOK, paginationResponse)
}

func (h *handlers) CreateTigerSightingHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request body to get the tiger sighting data
	if err := r.ParseMultipartForm(10 << 20); err != nil { // Max memory of 10 MB for file uploads
		utils.RespondWithError(w, http.StatusBadRequest, "Unable to parse form data")
		return
	}

	// Get the form values
	tigerIDStr := r.FormValue("tigerID")
	timestampStr := r.FormValue("timestamp")
	latStr := r.FormValue("lat")
	longStr := r.FormValue("long")

	// Convert the form values to appropriate types
	tigerID, err := strconv.Atoi(tigerIDStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid tigerID value")
		return
	}

	timestamp, err := time.Parse(time.RFC3339, timestampStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid timestamp value")
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid lat value")
		return
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid long value")
		return
	}

	reporterEmail, ok := auth.GetEmailFromContext(r.Context())
	if !ok {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve previous sighting")
		return
	}

	newSighting := models.TigerSighting{TigerID: tigerID, Timestamp: timestamp, Lat: lat, Long: long, ReporterEmail: reporterEmail}

	imageFile, _, err := r.FormFile("image")
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Failed to get image file")
		return
	}
	defer imageFile.Close()

	resizedImage, err := getProcessedImage(imageFile)
	if err != nil {
		h.Logger.Printf("Got Error Resizing Image: %v", err)
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	newSighting.Image = resizedImage
	err = h.TigerService.CreateTigerSightingService(&newSighting)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]interface{}{"message": "success"})
}

func getProcessedImage(imageFile multipart.File) ([]byte, error) {
	// Read the image data into a byte slice
	imageData, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return nil, err
	}

	// Resize the image to 250x200
	resizedImage, err := utils.ResizeImage(imageData, 250, 200)
	if err != nil {
		return nil, err
	}
	return resizedImage, nil
}

func (h *handlers) GetTigerSightingsByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tigerID := vars["id"]
	if tigerID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Missing tiger_id query parameter")
		return
	}

	// Convert the tiger ID to an integer
	tigerIDInt, err := strconv.Atoi(tigerID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid tiger_id query parameter")
		return
	}

	// Get the pagination parameters from the query string
	pageStr := r.FormValue("page")
	pageSizeStr := r.FormValue("pageSize")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = DefaultPageSize
	}

	tigerSightings, totalCount, err := h.TigerService.GetTigerSightingsByIDService(tigerIDInt, page, pageSize)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, t := range tigerSightings {
		img, _, err := image.Decode(bytes.NewReader(t.Image))
		if err != nil {
			fmt.Errorf("Error decoding image data: %v", err)
		}

		// Save the image to a new file
		fileName := fmt.Sprintf("%v_%v_%v_%v.jpeg", t.TigerID, t.Lat, t.Long, t.ReporterEmail)
		outputFile, err := os.Create(fileName) // we could have store it in S3 bucket, for simplicity storing it here.
		if err != nil {
			fmt.Errorf("Error creating output file: %v", err)
		}
		defer outputFile.Close()

		// Save the image in JPEG format
		if img != nil {
			err = jpeg.Encode(outputFile, img, nil)
			if err != nil {
				fmt.Errorf("Error encoding image data to file: %v", err)
			}
		}

		t.ImageFile = fileName
		t.Image = nil
	}

	// Construct pagination response
	paginationResponse := pagination{
		"page":           page,
		"pageSize":       pageSize,
		"totalCount":     totalCount,
		"totalPages":     int(math.Ceil(float64(totalCount) / float64(pageSize))),
		"tigerSightings": tigerSightings,
	}

	// Respond with the tiger sightings as JSON
	utils.RespondWithJSON(w, http.StatusOK, paginationResponse)
}
