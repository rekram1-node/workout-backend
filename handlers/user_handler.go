package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/rekram1-node/workout-backend/auth"
	"github.com/rekram1-node/workout-backend/models"
	"github.com/rekram1-node/workout-backend/repository"
	"github.com/rs/zerolog"
)

type LoginRepository interface {
	FindUserByCredentials(ctx context.Context, signInRequest repository.UserSignInRequest) (*models.User, error)
}

func LoginHandler(db LoginRepository, secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		var signinReq repository.UserSignInRequest
		if err := json.NewDecoder(r.Body).Decode(&signinReq); err != nil {
			logger.Warn().Err(err).Msg("failed to unmarshal body request")
			writeResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if err := validateRequest(signinReq); err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		user, err := db.FindUserByCredentials(ctx, signinReq)
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, map[string]string{"error": "unable to locate user"})
			return
		}

		token, err := auth.CreateAccessToken(user, secret)
		if err != nil {
			logger.Error().Err(err).Msg("failed to create access token")
			writeResponse(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}

		writeResponse(w, http.StatusOK, map[string]string{
			"token": token,
		})
	}
}

type UserRepository interface {
	CreateUser(ctx context.Context, userRequest repository.UserCreateRequest) (*models.User, error)
	ReadUser(ctx context.Context, uuid string) (*models.User, error)
	DeleteUser(ctx context.Context, uuid string) error
	UpdateUser(ctx context.Context, uuid string, updatedUser repository.UserUpdateRequest) error
}

func UserCreate(repo UserRepository, secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		var userReq repository.UserCreateRequest

		if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
			logger.Warn().Err(err).Msg("failed to unmarshal body request")
			writeResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if err := validateRequest(userReq); err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		user, err := repo.CreateUser(ctx, userReq)
		if err != nil {
			logger.Warn().Err(err).Msg("failed to create user")
			writeResponse(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}

		token, err := auth.CreateAccessToken(user, secret)
		if err != nil {
			logger.Error().Err(err).Msg("failed to create access token")
			writeResponse(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			return
		}

		writeResponse(w, http.StatusOK, map[string]string{
			"token":   token,
			"message": "successfully created user",
		})
	}
}

func UserUpdate(repo UserRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		userUUID := r.Header.Get("UUID")
		if userUUID == "" {
			logger.Error().Msg("unable to read user uuid from context")
			writeResponse(w, http.StatusInternalServerError, map[string]string{
				"error": "invalid token",
			})
			return
		}
		var updatedUser repository.UserUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
			logger.Debug().Err(err).Msg("unable to bind client struct")
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error": "failed to marshal to user",
			})
			return
		}
		if err := validateRequest(updatedUser); err != nil {
			writeResponse(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}
		if err := repo.UpdateUser(ctx, userUUID, updatedUser); err != nil {
			logger.Error().Err(err).Msgf("failed to update user: %s", userUUID)
			writeResponse(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to update user",
			})
			return
		}

		logger.Info().Msg("successfully updated user")
		writeResponse(w, http.StatusOK, map[string]string{
			"msg": "updated client",
		})
	}
}

func UserRead(repo UserRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		userUUID := r.Header.Get("UUID")
		user, err := repo.ReadUser(ctx, userUUID)
		if err != nil {
			logger.Error().Err(err).Msgf("failed to read user info for user: %s", userUUID)
		}

		writeResponse(w, http.StatusOK, *user)
	}
}

func UserDelete(repo UserRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		userUUID := r.Header.Get("UUID")
		if userUUID == "" {
			logger.Error().Msg("unable to read user uuid from context")
			writeResponse(w, http.StatusInternalServerError, map[string]string{
				"error": "invalid token",
			})
			return
		}

		if err := repo.DeleteUser(ctx, userUUID); err != nil {
			logger.Error().Err(err).Msgf("failed to delete user: %s", userUUID)
			writeResponse(w, http.StatusInternalServerError, map[string]string{
				"error": "failed to delete user: " + userUUID,
			})
			return
		}

		writeResponse(w, http.StatusOK, map[string]string{
			"message": "successfully deleted user: " + userUUID,
		})
	}
}
