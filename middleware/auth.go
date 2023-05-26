package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rekram1-node/workout-backend/auth"
	"github.com/rs/zerolog"
)

type JWTAuthentication struct {
	SecretKey string
}

func (jwtAuth JWTAuthentication) Authentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		ctx := r.Context()
		logger := zerolog.Ctx(ctx)
		t := strings.Split(authHeader, " ")
		if len(t) != 2 {
			logger.Debug().Msg("invalid token, missing bearer or token value")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid token",
			})
			return
		}

		authToken := t[1]
		authorized, err := auth.IsAuthorized(authToken, jwtAuth.SecretKey)
		if err != nil {
			logger.Info().Err(err).Msg("unathorized or failed to read token")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		if !authorized {
			logger.Info().Err(err).Msg("unathorized token")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "unathorized or invalid token",
			})
			return
		}

		uuid, err := auth.ReadUUIDFromToken(authToken, jwtAuth.SecretKey)
		if err != nil {
			logger.Info().Err(err).Msg("failed to read token")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": err.Error(),
			})
			return
		}

		if uuid == "" {
			logger.Info().Msg("missing uuid")
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "missing uuid in claims",
			})
			return
		}

		r.Header.Set("UUID", uuid)
		h.ServeHTTP(w, r)
	})
}
