package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/rekram1-node/workout-backend/models"
	"github.com/rekram1-node/workout-backend/repository"
	"github.com/rs/zerolog"
)

type MesoRepository interface {
	CreateMeso(ctx context.Context, mesoCreateReq *repository.MesoCreateRequest) (*models.Meso, error)
	ReadMeso(ctx context.Context, userUUID, mesoUUID string) (*repository.MesoResponse, error)
	ReadxMesos(ctx context.Context, userUUID string, mesoCount int) (*[]repository.MesoResponse, error)
	UpdateMeso(ctx context.Context, mesoUpdateReq *repository.MesoUpdateRequest) (*repository.MesoResponse, error)
	DeleteMeso(ctx context.Context, userUUID, mesoUUID string) error
}

func MesoCreate(repo MesoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		userUUID := r.Header.Get("UUID")
		var newMesoReq *repository.MesoCreateRequest

		err := json.NewDecoder(r.Body).Decode(&newMesoReq)
		switch {
		case errors.Is(err, io.EOF):
			writeResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request empty request body"})
			return
		case err != nil:
			logger.Warn().Err(err).Msg("failed to unmarshal body request")
			writeResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		newMesoReq.UserUUID = userUUID
		if err := validateRequest(newMesoReq); err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		meso, err := repo.CreateMeso(ctx, newMesoReq)
		if err != nil {
			logger.Warn().Err(err).Msg("failed to create user")
			writeResponse(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		logger.Info().Str("mesoUUID", meso.UUID).Msg("successfully created new meso")
		writeResponse(w, http.StatusOK, meso)
	}
}

func MesoRead(repo MesoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		userUUID := r.Header.Get("UUID")
		mesoUUID := r.URL.Query().Get("mesoUUID")
		if mesoUUID == "" {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error": "missing mesoUUID",
			})
			return
		}

		meso, err := repo.ReadMeso(ctx, userUUID, mesoUUID)
		if err != nil {
			logger.Error().Err(err).Msg("failed to find meso")
			writeResponse(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeResponse(w, http.StatusOK, meso)
	}
}

func MesosRead(repo MesoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		userUUID := r.Header.Get("UUID")
		numMesos, err := strconv.Atoi(r.URL.Query().Get("count"))
		if err != nil || numMesos <= 0 {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error": "invalid number of mesos",
			})
			return
		}
		mesos, err := repo.ReadxMesos(ctx, userUUID, numMesos)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to read top %v mesos", numMesos)
			writeResponse(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeResponse(w, http.StatusOK, mesos)
	}
}

func UpdateMeso(repo MesoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		userUUID := r.Header.Get("UUID")
		mesoUUID := r.URL.Query().Get("mesoUUID")
		if mesoUUID == "" {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error": "missing mesoUUID",
			})
			return
		}

		var newMesoReq *repository.MesoUpdateRequest

		err := json.NewDecoder(r.Body).Decode(&newMesoReq)
		switch {
		case errors.Is(err, io.EOF):
			writeResponse(w, http.StatusBadRequest, map[string]string{"error": "invalid request: empty request body"})
			return
		case err != nil:
			logger.Warn().Err(err).Msg("failed to unmarshal body request")
			writeResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		newMesoReq.UserUUID = userUUID
		newMesoReq.MesoUUID = mesoUUID
		if err := validateRequest(newMesoReq); err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		meso, err := repo.UpdateMeso(ctx, newMesoReq)
		if err != nil {
			logger.Error().Err(err).Msg("failed to update meso")
			writeResponse(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeResponse(w, http.StatusOK, meso)
	}
}

func DeleteMeso(repo MesoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		userUUID := r.Header.Get("UUID")
		mesoUUID := r.URL.Query().Get("mesoUUID")
		if mesoUUID == "" {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error": "missing mesoUUID",
			})
			return
		}

		if err := repo.DeleteMeso(ctx, userUUID, mesoUUID); err != nil {
			logger.Error().Err(err).Msg("failed to delete meso")
			writeResponse(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}

		writeResponse(w, http.StatusOK, map[string]string{
			"message": "successfully deleted meso: " + mesoUUID,
		})
	}
}
