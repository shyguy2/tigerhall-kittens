package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"tigerhall-kittens-app/pkg/models"
)

func TestGetMails(t *testing.T) {
	// Create a slice of previousSightings to be used for testing
	previousSightings := []*models.TigerSighting{
		{
			TigerID:       1,
			Lat:           40.7128,
			Long:          -74.0060,
			ReporterEmail: "test1@example.com",
		},
		{
			TigerID:       2,
			Lat:           34.0522,
			Long:          -118.2437,
			ReporterEmail: "test2@example.com",
		},
	}

	expectedEmails := []EmailTemplate{
		{
			Sub:       "Tiger Sights",
			Body:      "Tiger_1 is found at {Lat: 40.7128,Long: -74.006}",
			Recipient: "test1@example.com",
		},
		{
			Sub:       "Tiger Sights",
			Body:      "Tiger_2 is found at {Lat: 34.0522,Long: -118.2437}",
			Recipient: "test2@example.com",
		},
	}

	emailsJSON := GetMails(previousSightings)
	assert.NotNil(t, emailsJSON)

	// Unmarshal the JSON to EmailTemplate slice for comparison
	var actualEmails []EmailTemplate
	err := json.Unmarshal(emailsJSON, &actualEmails)
	assert.NoError(t, err)
	assert.Equal(t, expectedEmails, actualEmails)
}

func TestCalculateDistance(t *testing.T) {
	// Define coordinates for testing
	coord1 := models.Coordinates{Lat: 40.7128, Long: -74.0060}
	coord2 := models.Coordinates{Lat: 34.0522, Long: -118.2437}

	// Expected distance between the coordinates is approximately 3949.56 km
	expectedDistance := 3935.746254609722

	// Calculate the distance between the coordinates
	distance := CalculateDistance(coord1, coord2)
	assert.InDelta(t, expectedDistance, distance, 0.1)
}

//// Mock struct for image.Image
//type mockImage struct{}
//
//func (m *mockImage) ColorModel() color.Model {
//	return nil
//}
//
//func (m *mockImage) Bounds() image.Rectangle {
//	return image.Rectangle{}
//}
//
//func (m *mockImage) At(x, y int) color.Color {
//	return nil
//}
//
//func TestResizeImage(t *testing.T) {
//	// Encode a mock image into bytes
//	mockImgBytes := []byte{0, 1, 2, 3, 4, 5}
//
//	// Perform the image resizing
//	resizedImgBytes, err := ResizeImage(mockImgBytes, 250, 200)
//	assert.NoError(t, err)
//	assert.NotNil(t, resizedImgBytes)
//
//	// Test the decoding and resizing of the image
//	// Note: You can't fully validate the result without proper image decoding
//	// This test only checks if the function works without errors
//	_, _, err = image.Decode(bytes.NewReader(resizedImgBytes))
//	assert.NoError(t, err)
//}

func TestRespondWithError(t *testing.T) {
	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the RespondWithError function
	RespondWithError(rr, http.StatusNotFound, "Not Found")

	// Check if the response status code is 404 (Not Found)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Check if the response body contains the expected error message
	expectedResponse := map[string]string{"error": "Not Found"}
	var actualResponse map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &actualResponse)
	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, actualResponse)
}

func TestRespondWithJSON(t *testing.T) {
	// Create a test payload
	payload := map[string]interface{}{
		"message": "Hello, World!",
		"status":  true,
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the RespondWithJSON function
	RespondWithJSON(rr, http.StatusOK, payload)

	// Check if the response status code is 200 (OK)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check if the response body contains the expected JSON payload
	var responseBody map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &responseBody)
	assert.NoError(t, err)
	assert.Equal(t, payload, responseBody)
}
