package repository

import (
	"context"
	"fmt"

	"github.com/rekram1-node/workout-backend/models"
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	gormDB *gorm.DB
}

const (
	OPERATION = "operation"
	CREATE    = "create"
	READ      = "read"
	UPDATE    = "update"
	DELETE    = "delete"
	UUID      = "uuid"
)

func New(dbURI string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Meso{},
	); err != nil {
		return nil, err
	}

	return &Repository{
		gormDB: db,
	}, nil
}

func checkDBError(db *gorm.DB) error {
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected == 0 {
		return fmt.Errorf("zero rows affected in database operation")
	}

	return nil
}

func getDBLogger(repo *Repository, ctx context.Context, op, uuid string) (*gorm.DB, zerolog.Logger) {
	logger := zerolog.Ctx(ctx).With().
		Str(UUID, uuid).
		Str(OPERATION, op).
		Logger()
	database := repo.gormDB.WithContext(ctx)
	return database, logger
}
