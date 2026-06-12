package skill

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dionis-the-shark/apms-task-tracker/internal/models/skill"
	skillservice "github.com/dionis-the-shark/apms-task-tracker/internal/service/skill"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func Create(svc *skillservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sk skill.Skill
		if err := json.NewDecoder(r.Body).Decode(&sk); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := svc.Create(&sk); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(sk)
	}
}

func GetByID(svc *skillservice.Service) http.HandlerFunc {
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

func GetAll(svc *skillservice.Service) http.HandlerFunc {
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

func Update(svc *skillservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		skillID, err := uuid.Parse(chi.URLParam(r, "skill_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		var sk skill.Skill
		if err := json.NewDecoder(r.Body).Decode(&sk); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		sk.SkillID = skillID

		if err := svc.Update(sk); err != nil {
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

func Delete(svc *skillservice.Service) http.HandlerFunc {
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
