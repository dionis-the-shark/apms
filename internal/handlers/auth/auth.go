package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	authusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/auth/usecase"
)

func Register(svc authusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input authusecase.RegisterInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		resp, err := svc.Register(r.Context(), input)
		if err != nil {
			if errors.Is(err, authusecase.ErrEmailExists) || errors.Is(err, authusecase.ErrInvalidInput) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func Login(svc authusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input authusecase.LoginInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		resp, err := svc.Login(r.Context(), input)
		if err != nil {
			if errors.Is(err, authusecase.ErrInvalidCredentials) || errors.Is(err, authusecase.ErrInvalidInput) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func Validate(svc authusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := extractBearerToken(r.Header.Get("Authorization"))
		if token == "" {
			http.Error(w, "Missing or invalid Authorization header", http.StatusBadRequest)
			return
		}

		resp, err := svc.ValidateToken(r.Context(), token)
		if err != nil {
			if errors.Is(err, authusecase.ErrInvalidToken) || errors.Is(err, authusecase.ErrTokenExpired) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}
}

func extractBearerToken(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
