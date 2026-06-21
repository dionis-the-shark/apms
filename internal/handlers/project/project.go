package project

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dionis-the-shark/apms-task-tracker/internal/authz"
	"github.com/dionis-the-shark/apms-task-tracker/internal/httpauth"
	authusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/auth/usecase"
	projectusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/project/usecase"
	roleusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func Create(svc projectusecase.API, authSvc authusecase.API, roleSvc roleusecase.API) http.HandlerFunc {
	return httpauth.RequirePermission(authSvc, roleSvc, authz.ProjectCreate, func(w http.ResponseWriter, r *http.Request) {
		var input projectusecase.CreateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		p, err := svc.Create(ctx, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(p)
	})
}

func GetByID(svc projectusecase.API, authSvc authusecase.API) http.HandlerFunc {
	return httpauth.RequireAuth(authSvc, func(w http.ResponseWriter, r *http.Request) {
		projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		p, err := svc.GetByID(ctx, projectID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(p)
	})
}

func GetAll(svc projectusecase.API, authSvc authusecase.API) http.HandlerFunc {
	return httpauth.RequireAuth(authSvc, func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		projects, err := svc.GetAll(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(projects)
	})
}

func Update(svc projectusecase.API, authSvc authusecase.API) http.HandlerFunc {
	return httpauth.RequireAuth(authSvc, func(w http.ResponseWriter, r *http.Request) {
		projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		var input projectusecase.UpdateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		p, err := svc.Update(ctx, projectID, input)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(p)
	})
}

func Delete(svc projectusecase.API, authSvc authusecase.API) http.HandlerFunc {
	return httpauth.RequireAuth(authSvc, func(w http.ResponseWriter, r *http.Request) {
		projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		if err := svc.Delete(ctx, projectID); err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
