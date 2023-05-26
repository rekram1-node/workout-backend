package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rekram1-node/workout-backend/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserSignInRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (repo *Repository) FindUserByCredentials(ctx context.Context, signInRequest UserSignInRequest) (*models.User, error) {
	logger := zerolog.Ctx(ctx).With().Str("user", signInRequest.Username).Logger()
	user := &models.User{
		Username: signInRequest.Username,
		Password: signInRequest.Password,
	}
	res := repo.gormDB.WithContext(ctx).Where("username = ? AND password = ?", signInRequest.Username, signInRequest.Password).
		Find(&user)
	if err := checkDBError(res); err != nil {
		logger.Error().Err(err).Msg("unable to find user")
		return nil, err
	}

	return user, nil
}

type UserCreateRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (repo *Repository) CreateUser(ctx context.Context, userRequest UserCreateRequest) (*models.User, error) {
	user := &models.User{
		UUID:     uuid.New().String(),
		Username: userRequest.Username,
		Password: userRequest.Password,
	}
	gormDB, logger := getDBLogger(repo, ctx, CREATE, user.UUID)
	dberr := gormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, obj := range []interface{}{user} {
			if err := checkDBError(tx.WithContext(ctx).Debug().Create(obj)); err != nil {
				logger.Error().Err(err).Msg("failed to create user")
				return err
			}
		}

		return nil
	})

	if errors.Is(dberr, gorm.ErrDuplicatedKey) {
		return user, fmt.Errorf("username already exists")
	}
	if dberr != nil {
		logger.Error().Err(dberr).Msg("database error creating new user")
		return user, dberr
	}

	logger.Info().Interface("user", user).Msg("created new user")
	return user, dberr
}

func (repo *Repository) ReadUser(ctx context.Context, uuid string) (*models.User, error) {
	var user *models.User
	gormDB, logger := getDBLogger(repo, ctx, READ, uuid)

	res := gormDB.WithContext(ctx).Where("uuid = ?", uuid).
		Find(&user)
	if err := checkDBError(res); err != nil {
		logger.Error().Err(err).Msg("unable to find user")
		return nil, err
	}

	return user, nil
}

func (repo *Repository) DeleteUser(ctx context.Context, uuid string) error {
	gormDB, logger := getDBLogger(repo, ctx, DELETE, uuid)
	dberr := gormDB.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		var user *models.User
		var mesoIDs []uint

		dbUsr := tx.Where("uuid = ?", uuid).
			Take(&user)
		if err := checkDBError(dbUsr); err != nil {
			logger.Error().Err(err).Msg("unable to lookup user")
			return err
		}

		dbMesoResult := tx.Model(&models.Meso{}).
			Where("user_id = ?", user.ID).
			Pluck("id", &mesoIDs)
		if err := checkDBError(dbMesoResult); err != nil {
			logger.Error().Err(err).Msg("unable to lookup user meso uuids")
			return err
		}

		for _, mesoID := range mesoIDs {
			resultDeleteMeso := tx.Model(&models.Meso{}).
				Where("id = ?", mesoID).
				Delete(&models.Meso{})
			if err := checkDBError(resultDeleteMeso); err != nil {
				logger.Error().Err(err).Uint("meso_id", mesoID).Msg("database error deleting meso")
				return err
			}
		}

		resultDelete := tx.
			Select("Mesos").
			Where("id = ?", user.ID).
			Delete(&user)
		if err := checkDBError(resultDelete); err != nil {
			logger.Error().Err(err).Msg("database error deleting user")
			return err
		}

		return nil
	})

	if dberr != nil {
		logger.Error().Err(dberr).Msg("database transaction failure deleting user")
		return dberr
	}

	return nil
}

type UserUpdateRequest struct {
	UUID     uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Username string    `json:"username" gorm:"uniqueIndex;not null" validate:"required"`
	Password string    `json:"password" validate:"required"`
}

func (repo *Repository) UpdateUser(ctx context.Context, uuid string, userUpdate UserUpdateRequest) error {
	var user *models.User
	gormDB, logger := getDBLogger(repo, ctx, UPDATE, uuid)

	dberr := gormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := make(map[string]interface{})
		if userUpdate.Username != "" {
			updates["username"] = userUpdate.Username
		}

		res := tx.WithContext(ctx).
			Model(&user).
			Where("uuid = ?", userUpdate.UUID).
			Updates(updates)
		if err := checkDBError(res); err != nil {
			logger.Error().Err(err).Msg("error updating user metadata")
			return err
		}

		return nil
	})

	if dberr != nil {
		logger.Error().Err(dberr).Msg("error updating user in database")
		return dberr
	}

	return nil
}
