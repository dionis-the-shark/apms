package developer_skill

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	developerskillusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func Create(svc developerskillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input developerskillusecase.CreateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		ds, err := svc.Create(ctx, input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(ds)
	}
}

func GetByIDs(svc developerskillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		developerID, err := uuid.Parse(chi.URLParam(r, "developer_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}
		skillID, err := uuid.Parse(chi.URLParam(r, "skill_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		ds, err := svc.GetByIDs(ctx, developerID, skillID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ds)
	}
}

func GetAll(svc developerskillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		skills, err := svc.GetAll(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(skills)
	}
}

func GetByDeveloper(svc developerskillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		developerID, err := uuid.Parse(chi.URLParam(r, "developer_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		skills, err := svc.GetByDeveloper(ctx, developerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(skills)
	}
}

func Update(svc developerskillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		oldDeveloperID, err := uuid.Parse(chi.URLParam(r, "developer_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}
		oldSkillID, err := uuid.Parse(chi.URLParam(r, "skill_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		var input developerskillusecase.UpdateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		ds, err := svc.Update(ctx, oldDeveloperID, oldSkillID, input)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(ds)
	}
}

func Delete(svc developerskillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		developerID, err := uuid.Parse(chi.URLParam(r, "developer_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}
		skillID, err := uuid.Parse(chi.URLParam(r, "skill_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		if err := svc.Delete(ctx, developerID, skillID); err != nil {
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
