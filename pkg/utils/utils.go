package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"net/http"

	"github.com/disintegration/imaging"
	"github.com/umahmood/haversine"
	"tigerhall-kittens-app/pkg/models"
)

type EmailTemplate struct {
	Sub       string `json:"subject"`
	Body      string `json:"body"`
	Recipient string `json:"recipient"`
}

func GetMails(previousSightings []*models.TigerSighting) []byte {
	var emails []EmailTemplate
	for _, pr := range previousSightings {
		emails = append(emails, EmailTemplate{
			Sub:       "Tiger Sights",
			Body:      fmt.Sprintf(`Tiger_%v is found at {Lat: %v,Long: %v}`, pr.TigerID, pr.Lat, pr.Long),
			Recipient: pr.ReporterEmail,
		})
	}

	emailsJSON, err := json.Marshal(emails)
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return nil
	}
	return []byte(emailsJSON)
}

func CalculateDistance(coord1, coord2 models.Coordinates) float64 {
	point1 := haversine.Coord{Lat: coord1.Lat, Lon: coord1.Long}
	point2 := haversine.Coord{Lat: coord2.Lat, Lon: coord2.Long}

	_, distanceKm := haversine.Distance(point1, point2)

	return distanceKm
}

func ResizeImage(imageBytes []byte, width, height int) ([]byte, error) {
	// Decode the imageBytes into an image.Image
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, err
	}

	// Resize the image using the Lanczos filter
	resizedImg := imaging.Resize(img, width, height, imaging.Lanczos)

	// Encode the resized image back to bytes
	var buf bytes.Buffer
	err = imaging.Encode(&buf, resizedImg, imaging.JPEG)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	response := map[string]string{"error": message}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	jsonResponse, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonResponse)
}
