package task_dependency

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	taskdependencyusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task_dependency/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func Create(svc taskdependencyusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input taskdependencyusecase.CreateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		td, err := svc.Create(ctx, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(td)
	}
}

func GetByIDs(svc taskdependencyusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blockedTaskID, err := uuid.Parse(chi.URLParam(r, "blocked_task_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}
		dependentOnID, err := uuid.Parse(chi.URLParam(r, "dependent_on_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		td, err := svc.GetByIDs(ctx, blockedTaskID, dependentOnID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(td)
	}
}

func GetAll(svc taskdependencyusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		deps, err := svc.GetAll(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(deps)
	}
}

func GetByBlockedTask(svc taskdependencyusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blockedTaskID, err := uuid.Parse(chi.URLParam(r, "blocked_task_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		deps, err := svc.GetByBlockedTask(ctx, blockedTaskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(deps)
	}
}

func Update(svc taskdependencyusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oldBlockedTaskID, err := uuid.Parse(chi.URLParam(r, "blocked_task_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}
		oldDependentOnID, err := uuid.Parse(chi.URLParam(r, "dependent_on_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		var input taskdependencyusecase.UpdateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		td, err := svc.Update(ctx, oldBlockedTaskID, oldDependentOnID, input)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(td)
	}
}

func Delete(svc taskdependencyusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blockedTaskID, err := uuid.Parse(chi.URLParam(r, "blocked_task_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}
		dependentOnID, err := uuid.Parse(chi.URLParam(r, "dependent_on_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		if err := svc.Delete(ctx, blockedTaskID, dependentOnID); err != nil {
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
