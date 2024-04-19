package repository

import (
	"tigerhall-kittens-app/pkg/models"
	"tigerhall-kittens-app/pkg/repository/store"
)

type TigerRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	CreateTiger(tiger *models.Tiger) error
	GetAllTigersWithPagination(page, pageSize int) ([]*models.Tiger, int, error)
	CreateTigerSighting(tigerSighting *models.TigerSighting) error
	GetTigerSightingsByID(tigerID int) ([]*models.TigerSighting, error)
	GetPreviousTigerSighting(tigerID int) (*models.TigerSighting, error)
	GetTigerSightingsByIDWithPagination(tigerID, page, pageSize int) ([]*models.TigerSighting, int, error)
}

func NewPostgresRepository(connection string) (TigerRepository, error) {
	db, err := store.NewPostgresDB(connection)
	return store.NewPostgresRepository(db), err
}
