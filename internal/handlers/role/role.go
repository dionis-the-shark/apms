package role

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	roleusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func Create(svc roleusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input roleusecase.CreateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		roleModel, err := svc.Create(ctx, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(roleModel)
	}
}

func GetByID(svc roleusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		roleModel, err := svc.GetByID(ctx, roleID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(roleModel)
	}
}

func GetAll(svc roleusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		roles, err := svc.GetAll(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(roles)
	}
}

func Update(svc roleusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		var input roleusecase.UpdateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		roleModel, err := svc.Update(ctx, roleID, input)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(roleModel)
	}
}

func Delete(svc roleusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleID, err := uuid.Parse(chi.URLParam(r, "role_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		if err := svc.Delete(ctx, roleID); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
