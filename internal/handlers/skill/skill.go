package skill

import (
	"database/sql"
	"encoding/json"
	"net/http"

	skillusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/skill/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func Create(svc skillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input skillusecase.CreateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		sk, err := svc.Create(input)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(sk)
	}
}

func GetByID(svc skillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		skillID, err := uuid.Parse(chi.URLParam(r, "skill_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		sk, err := svc.GetByID(skillID)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sk)
	}
}

func GetAll(svc skillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		skills, err := svc.GetAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(skills)
	}
}

func Update(svc skillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		skillID, err := uuid.Parse(chi.URLParam(r, "skill_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		var input skillusecase.UpdateInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		sk, err := svc.Update(skillID, input)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(sk)
	}
}

func Delete(svc skillusecase.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		skillID, err := uuid.Parse(chi.URLParam(r, "skill_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		if err := svc.Delete(skillID); err != nil {
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
