package service

import (
	"errors"
	"fmt"
	"log"
	"sort"

	"tigerhall-kittens-app/pkg/auth"
	"tigerhall-kittens-app/pkg/messaging"
	"tigerhall-kittens-app/pkg/models"
	"tigerhall-kittens-app/pkg/repository"
	"tigerhall-kittens-app/pkg/utils"
)

type service struct {
	TigerRepo     repository.TigerRepository
	messageBroker *messaging.MessageBroker
}

func NewTigerService(tigerRepository repository.TigerRepository, broker *messaging.MessageBroker) TigerService {
	return service{
		TigerRepo:     tigerRepository,
		messageBroker: broker,
	}
}

type TigerService interface {
	SignupService(*models.User) error
	LoginService(models.LoginCredentials) (*models.User, error)
	CreateTigerService(tiger models.Tiger) error
	GetAllTigersService(page, size int) ([]*models.Tiger, int, error)
	CreateTigerSightingService(*models.TigerSighting) error
	GetTigerSightingsByIDService(tigerID, page, pageSize int) ([]*models.TigerSighting, int, error)
}

func (s service) SignupService(user *models.User) error {
	// Hash the user's password before saving to the database
	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	user.Password = hashedPassword

	// Create the user in the database
	if err := s.TigerRepo.CreateUser(user); err != nil {
		return errors.New("failed to create user")
	}
	return err
}

func (s service) LoginService(credentials models.LoginCredentials) (*models.User, error) {
	// Find the user by email in the database
	user, err := s.TigerRepo.GetUserByEmail(credentials.Email)
	if err != nil {
		return &models.User{}, errors.New("invalid email or password")
	}

	// Verify the password
	if err := auth.VerifyPassword(user.Password, credentials.Password); err != nil {
		return &models.User{}, errors.New("invalid email or password")
	}
	return user, nil
}

func (s service) CreateTigerService(tiger models.Tiger) error {
	// Create the tiger in the database
	if err := s.TigerRepo.CreateTiger(&tiger); err != nil {
		return errors.New("failed to create tiger")
	}
	return nil
}

func (s service) GetAllTigersService(page, size int) ([]*models.Tiger, int, error) {
	// Get a list of all tigers from the database with pagination
	tigers, totalCount, err := s.TigerRepo.GetAllTigersWithPagination(page, size)
	if err != nil {
		return []*models.Tiger{}, totalCount, errors.New("failed to fetch tigers")
	}

	// Sort the tigers by the last seen time
	sort.Slice(tigers, func(i, j int) bool { return tigers[i].LastSeen.After(tigers[j].LastSeen) })
	return tigers, totalCount, nil
}

func (s service) CreateTigerSightingService(newSighting *models.TigerSighting) error {
	// Check if the required fields are provided
	if newSighting.Lat == 0 || newSighting.Long == 0 || newSighting.Timestamp.IsZero() || newSighting.ReporterEmail == "" {
		return errors.New("latitude, longitude, timestamp and reporterEmail are required")
	}

	// Check if the tiger has a previous sighting
	previousSighting, err := s.TigerRepo.GetPreviousTigerSighting(newSighting.TigerID)
	if err != nil {
		return errors.New("failed to retrieve previous sighting")
	}

	// If there is a previous sighting, calculate the distance between the coordinates
	if previousSighting != nil {
		previousCoordinates := models.Coordinates{Lat: previousSighting.Lat, Long: previousSighting.Long}
		currentCoordinates := models.Coordinates{Lat: newSighting.Lat, Long: newSighting.Long}
		distance := utils.CalculateDistance(previousCoordinates, currentCoordinates)

		// If the distance is less than or equal to 5 kilometers, reject the new sighting
		if distance <= 5.0 {
			return errors.New("A tiger sighting within 5 kilometers already exists")
		}
	}

	// Create the tiger sighting in the database
	err = s.TigerRepo.CreateTigerSighting(newSighting)
	if err != nil {
		return errors.New("failed to create tiger sighting")
	}

	previousSightings, err := s.TigerRepo.GetTigerSightingsByID(newSighting.TigerID)
	if err != nil {
		return errors.New("failed to retrieve previous sightings")
	}

	// Publish a new tiger sighting message
	if s.messageBroker != nil {
		if err := s.messageBroker.PublishMessage(utils.GetMails(previousSightings)); err != nil {
			log.Printf("failed to publish message: %v", err)
		}
	}

	return nil
}

func (s service) GetTigerSightingsByIDService(tigerID, page, pageSize int) ([]*models.TigerSighting, int, error) {
	// Get a list of all tiger sightings for the specific tiger from the database with pagination
	tigerSightings, totalCount, err := s.TigerRepo.GetTigerSightingsByIDWithPagination(tigerID, page, pageSize)
	if err != nil {
		return []*models.TigerSighting{}, totalCount, fmt.Errorf("error: %v", err)
	}

	// Sort the tiger sightings by date
	sort.Slice(tigerSightings, func(i, j int) bool {
		return tigerSightings[i].Timestamp.After(tigerSightings[j].Timestamp)
	})

	return tigerSightings, totalCount, nil
}
