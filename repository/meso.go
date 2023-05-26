package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rekram1-node/workout-backend/models"
	"gorm.io/gorm"
)

type MesoCreateRequest struct {
	UserUUID                                                       string `validate:"required"`
	MesoUUID                                                       string
	Name                                                           string      `validate:"required"`
	Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday *models.Day `validate:"required"`
}

func (repo *Repository) CreateMeso(ctx context.Context, mesoCreateReq *MesoCreateRequest) (*models.Meso, error) {
	gormDB, logger := getDBLogger(repo, ctx, CREATE, mesoCreateReq.UserUUID)
	var user *models.User
	var meso *models.Meso

	dberr := gormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.WithContext(ctx).Where("uuid = ?", mesoCreateReq.UserUUID).
			Preload("Mesos").Find(&user)

		if err := checkDBError(res); err != nil {
			logger.Error().Err(err).Msg("error finding new meso parent user")
			return err
		}

		mesoUUID := uuid.NewString()
		meso = &models.Meso{
			User:     user,
			UserUUID: user.UUID,
			UUID:     mesoUUID,
			Name:     mesoCreateReq.Name,
			Weeks: &[]models.Week{
				{
					Monday:    mesoCreateReq.Monday,
					Tuesday:   mesoCreateReq.Tuesday,
					Wednesday: mesoCreateReq.Wednesday,
					Thursday:  mesoCreateReq.Thursday,
					Friday:    mesoCreateReq.Friday,
					Saturday:  mesoCreateReq.Saturday,
					Sunday:    mesoCreateReq.Sunday,
				},
			},
		}

		for _, obj := range []interface{}{meso} {
			if err := checkDBError(tx.WithContext(ctx).Debug().Create(obj)); err != nil {
				logger.Error().Err(err).Msg("failed to create meso")
				return err
			}
		}

		return nil
	})

	if dberr != nil {
		logger.Error().Err(dberr).Msg("database error creating new meso")
		return nil, dberr
	}

	logger.Info().Str("meso_uuid", meso.UUID).Msg("created new meso")

	return meso, nil
}

type MesoResponse struct {
	Name  string
	UUID  string
	Weeks *[]models.Week
}

func (repo *Repository) ReadMeso(ctx context.Context, userUUID, mesoUUID string) (*MesoResponse, error) {
	gormDB, logger := getDBLogger(repo, ctx, READ, userUUID)
	logger = logger.With().Str("meso_uuid", mesoUUID).Logger()
	var meso *models.Meso
	res := gormDB.WithContext(ctx).
		Where("user_uuid = ? AND uuid = ?", userUUID, mesoUUID).
		Where("uuid = ?", mesoUUID).Find(&meso)
	if err := checkDBError(res); err != nil {
		logger.Error().Err(err).Msg("failed to find meso")
		return nil, err
	}

	mesoResponse := MesoResponse{
		Name:  meso.Name,
		UUID:  meso.UUID,
		Weeks: meso.Weeks,
	}

	return &mesoResponse, nil
}

func (repo *Repository) ReadxMesos(ctx context.Context, userUUID string, mesoCount int) (*[]MesoResponse, error) {
	gormDB, logger := getDBLogger(repo, ctx, READ, userUUID)
	var mesos []models.Meso
	foundMesos := []MesoResponse{}
	res := gormDB.WithContext(ctx).Where("user_uuid = ?", userUUID).
		Order("updated_at DESC").
		Limit(mesoCount).
		Find(&mesos)
	if err := checkDBError(res); err != nil {
		logger.Error().Err(err).Msg("error finding user")
		return nil, err
	}

	for _, meso := range mesos {
		mesoRes := MesoResponse{
			Name:  meso.Name,
			UUID:  meso.UUID,
			Weeks: meso.Weeks,
		}
		foundMesos = append(foundMesos, mesoRes)
	}

	return &foundMesos, nil
}

type MesoUpdateRequest struct {
	UserUUID string
	MesoUUID string
	Name     string
	Weeks    *[]models.Week
}

func (repo *Repository) UpdateMeso(ctx context.Context, mesoUpdateReq *MesoUpdateRequest) (*MesoResponse, error) {
	gormDB, logger := getDBLogger(repo, ctx, UPDATE, mesoUpdateReq.UserUUID)
	logger = logger.With().Str("meso_uuid", mesoUpdateReq.MesoUUID).Logger()
	var meso models.Meso
	res := gormDB.Where("user_uuid = ?", mesoUpdateReq.UserUUID).Where("uuid = ?", mesoUpdateReq.MesoUUID).First(&meso)
	if res.Error != nil || res.RowsAffected == 0 {
		return nil, fmt.Errorf("no meso found with uuid [%s] for user [%s]: %w", mesoUpdateReq.MesoUUID, mesoUpdateReq.UserUUID, res.Error)
	}

	dberr := repo.gormDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		updates := make(map[string]interface{})
		if mesoUpdateReq.Name != "" {
			updates["name"] = mesoUpdateReq.Name
		}
		if mesoUpdateReq.Weeks != nil {
			updates["week"] = mesoUpdateReq.Weeks
		}

		res := tx.WithContext(ctx).Model(&meso).
			Where("user_uuid = ?", mesoUpdateReq.UserUUID).
			Where("uuid = ?", mesoUpdateReq.MesoUUID).
			Updates(updates).Find(&meso)

		if err := checkDBError(res); err != nil {
			logger.Error().Err(err).Msg("error updating meso metadata")
			return err
		}

		return nil
	})

	if dberr != nil {
		logger.Error().Err(dberr).Msg("error updating meso")
		return nil, dberr
	}

	return repo.ReadMeso(ctx, mesoUpdateReq.UserUUID, mesoUpdateReq.MesoUUID)
}

func (repo *Repository) DeleteMeso(ctx context.Context, userUUID, mesoUUID string) error {
	db, logger := getDBLogger(repo, ctx, DELETE, mesoUUID)
	logger = logger.With().Str("user", userUUID).Logger()
	dberr := db.Transaction(func(tx *gorm.DB) error {
		tx = tx.WithContext(ctx)
		var meso models.Meso

		dbMeso := tx.
			Where("user_uuid", userUUID).
			Where("uuid", mesoUUID).
			First(&meso)
		if err := checkDBError(dbMeso); err != nil {
			logger.Error().Err(err).Msg("unable to lookup meso")
			return err
		}

		resultDelete := tx.
			Where("uuid = ?", mesoUUID).
			Delete(&meso)
		if err := checkDBError(resultDelete); err != nil {
			logger.Error().Err(err).Msg("database error deleting meso")
			return err
		}

		return nil
	})
	if dberr != nil {
		logger.Error().Err(dberr).Msg("database transaction error deleting meso")
		return dberr
	}

	return nil
}
