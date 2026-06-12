package developer_skill

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/dionis-the-shark/apms-task-tracker/internal/models/developer_skill"
	developerskillservice "github.com/dionis-the-shark/apms-task-tracker/internal/service/developer_skill"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func Create(svc *developerskillservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var ds developer_skill.DeveloperSkill
		if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := svc.Create(ds); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(ds)
	}
}

func GetByIDs(svc *developerskillservice.Service) http.HandlerFunc {
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

		ds, err := svc.GetByIDs(developerID, skillID)
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

func GetAll(svc *developerskillservice.Service) http.HandlerFunc {
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

func GetByDeveloper(svc *developerskillservice.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		developerID, err := uuid.Parse(chi.URLParam(r, "developer_id"))
		if err != nil {
			http.Error(w, "Invalid UUID", http.StatusBadRequest)
			return
		}

		skills, err := svc.GetByDeveloper(developerID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(skills)
	}
}

func Update(svc *developerskillservice.Service) http.HandlerFunc {
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

		var ds developer_skill.DeveloperSkill
		if err := json.NewDecoder(r.Body).Decode(&ds); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := svc.Update(oldDeveloperID, oldSkillID, ds.DeveloperID, ds.SkillID); err != nil {
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

func Delete(svc *developerskillservice.Service) http.HandlerFunc {
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

		if err := svc.Delete(developerID, skillID); err != nil {
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
