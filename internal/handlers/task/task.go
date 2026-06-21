package task

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/dionis-the-shark/apms-task-tracker/internal/authz"
	"github.com/dionis-the-shark/apms-task-tracker/internal/httpauth"
	authusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/auth/usecase"
	roleusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role/usecase"
	taskusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func Create(svc taskusecase.API, authSvc authusecase.API, roleSvc roleusecase.API) http.HandlerFunc {
	return httpauth.RequirePermission(authSvc, roleSvc, authz.TaskCreate, func(w http.ResponseWriter, r *http.Request) {
		var input taskusecase.CreateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		t, err := svc.Create(ctx, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(t)
	})
}

func GetByID(svc taskusecase.API, authSvc authusecase.API, roleSvc roleusecase.API) http.HandlerFunc {
	return httpauth.RequirePermission(authSvc, roleSvc, authz.TaskRead, func(w http.ResponseWriter, r *http.Request) {
		taskID, err := uuid.Parse(chi.URLParam(r, "task_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		t, err := svc.GetByID(ctx, taskID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(t)
	})
}

func GetByProject(svc taskusecase.API, authSvc authusecase.API, roleSvc roleusecase.API) http.HandlerFunc {
	return httpauth.RequirePermission(authSvc, roleSvc, authz.TaskRead, func(w http.ResponseWriter, r *http.Request) {
		projectID, err := uuid.Parse(chi.URLParam(r, "project_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		tasks, err := svc.GetByProject(ctx, projectID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tasks)
	})
}

func Update(svc taskusecase.API, authSvc authusecase.API, roleSvc roleusecase.API) http.HandlerFunc {
	return httpauth.RequirePermission(authSvc, roleSvc, authz.TaskUpdate, func(w http.ResponseWriter, r *http.Request) {
		taskID, err := uuid.Parse(chi.URLParam(r, "task_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		var input taskusecase.UpdateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		t, err := svc.Update(ctx, taskID, input)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(t)
	})
}

func Delete(svc taskusecase.API, authSvc authusecase.API, roleSvc roleusecase.API) http.HandlerFunc {
	return httpauth.RequirePermission(authSvc, roleSvc, authz.TaskDelete, func(w http.ResponseWriter, r *http.Request) {
		taskID, err := uuid.Parse(chi.URLParam(r, "task_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		if err := svc.Delete(ctx, taskID); err != nil {
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
