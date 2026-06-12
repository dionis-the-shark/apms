package task_dependency

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dionis-the-shark/apms-task-tracker/internal/models/taskdependency"
	taskdependencyservice "github.com/dionis-the-shark/apms-task-tracker/internal/service/task_dependency"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func Create(svc *taskdependencyservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var td taskdependency.TaskDependency
		if err := json.NewDecoder(r.Body).Decode(&td); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := svc.Create(td); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(td)
	}
}

func GetByIDs(svc *taskdependencyservice.Service) http.HandlerFunc {
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

		td, err := svc.GetByIDs(blockedTaskID, dependentOnID)
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

func GetAll(svc *taskdependencyservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		deps, err := svc.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(deps)
	}
}

func GetByBlockedTask(svc *taskdependencyservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		blockedTaskID, err := uuid.Parse(chi.URLParam(r, "blocked_task_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		deps, err := svc.GetByBlockedTask(blockedTaskID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(deps)
	}
}

func Update(svc *taskdependencyservice.Service) http.HandlerFunc {
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

		var td taskdependency.TaskDependency
		if err := json.NewDecoder(r.Body).Decode(&td); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := svc.Update(oldBlockedTaskID, oldDependentOnID, td.BlockedTaskID, td.DependentOnID); err != nil {
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

func Delete(svc *taskdependencyservice.Service) http.HandlerFunc {
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

		if err := svc.Delete(blockedTaskID, dependentOnID); err != nil {
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
