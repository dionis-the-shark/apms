package http

import (
	nethttp "net/http"

	"github.com/dionis-the-shark/apms-task-tracker/internal/app"
	developer_skill_handlers "github.com/dionis-the-shark/apms-task-tracker/internal/handlers/developer_skill"
	project_handlers "github.com/dionis-the-shark/apms-task-tracker/internal/handlers/project"
	skill_handlers "github.com/dionis-the-shark/apms-task-tracker/internal/handlers/skill"
	task_handlers "github.com/dionis-the-shark/apms-task-tracker/internal/handlers/task"
	task_dependency_handlers "github.com/dionis-the-shark/apms-task-tracker/internal/handlers/task_dependency"
	user_handlers "github.com/dionis-the-shark/apms-task-tracker/internal/handlers/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(services *app.Services) nethttp.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", user_handlers.Create(services.Users))
		r.Get("/", user_handlers.GetAll(services.Users))
		r.Get("/{user_id}", user_handlers.GetByID(services.Users))
		r.Put("/{user_id}", user_handlers.Update(services.Users))
		r.Delete("/{user_id}", user_handlers.Delete(services.Users))
	})

	r.Route("/projects", func(r chi.Router) {
		r.Post("/", project_handlers.Create(services.Projects))
		r.Get("/", project_handlers.GetAll(services.Projects))
		r.Get("/{project_id}", project_handlers.GetByID(services.Projects))
		r.Put("/{project_id}", project_handlers.Update(services.Projects))
		r.Delete("/{project_id}", project_handlers.Delete(services.Projects))
		r.Get("/{project_id}/tasks", task_handlers.GetByProject(services.Tasks))
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", task_handlers.Create(services.Tasks))
		r.Get("/{task_id}", task_handlers.GetByID(services.Tasks))
		r.Put("/{task_id}", task_handlers.Update(services.Tasks))
		r.Delete("/{task_id}", task_handlers.Delete(services.Tasks))
	})

	r.Route("/skills", func(r chi.Router) {
		r.Post("/", skill_handlers.Create(services.Skills))
		r.Get("/", skill_handlers.GetAll(services.Skills))
		r.Get("/{skill_id}", skill_handlers.GetByID(services.Skills))
		r.Put("/{skill_id}", skill_handlers.Update(services.Skills))
		r.Delete("/{skill_id}", skill_handlers.Delete(services.Skills))
	})

	r.Route("/developer-skills", func(r chi.Router) {
		r.Post("/", developer_skill_handlers.Create(services.DeveloperSkills))
		r.Get("/", developer_skill_handlers.GetAll(services.DeveloperSkills))
		r.Get("/developers/{developer_id}", developer_skill_handlers.GetByDeveloper(services.DeveloperSkills))
		r.Get("/{developer_id}/{skill_id}", developer_skill_handlers.GetByIDs(services.DeveloperSkills))
		r.Put("/{developer_id}/{skill_id}", developer_skill_handlers.Update(services.DeveloperSkills))
		r.Delete("/{developer_id}/{skill_id}", developer_skill_handlers.Delete(services.DeveloperSkills))
	})

	r.Route("/task-dependencies", func(r chi.Router) {
		r.Post("/", task_dependency_handlers.Create(services.TaskDependencies))
		r.Get("/", task_dependency_handlers.GetAll(services.TaskDependencies))
		r.Get("/blocked/{blocked_task_id}", task_dependency_handlers.GetByBlockedTask(services.TaskDependencies))
		r.Get("/{blocked_task_id}/{dependent_on_id}", task_dependency_handlers.GetByIDs(services.TaskDependencies))
		r.Put("/{blocked_task_id}/{dependent_on_id}", task_dependency_handlers.Update(services.TaskDependencies))
		r.Delete("/{blocked_task_id}/{dependent_on_id}", task_dependency_handlers.Delete(services.TaskDependencies))
	})

	return r
}
