package models

import (
	"time"
)

type TigerSighting struct {
	ID            int       `json:"id"`
	TigerID       int       `json:"tigerID"`
	Timestamp     time.Time `json:"timestamp"`
	Lat           float64   `json:"lat"`
	Long          float64   `json:"long"`
	Image         []byte    `json:"image,omitempty"`
	ImageFile     string    `json:"imageFile,omitempty"`
	ReporterEmail string    `json:"reporterEmail"`
}
