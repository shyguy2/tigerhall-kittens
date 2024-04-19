package models

import "time"

type Tiger struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	DateOfBirth time.Time `json:"date_of_birth"`
	LastSeen    time.Time `json:"last_seen"`
	Lat         float64   `json:"lat"`
	Long        float64   `json:"long"`
}

type Coordinates struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}
