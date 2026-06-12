package app

import (
	"database/sql"

	"github.com/dionis-the-shark/apms-task-tracker/internal/repository"
	authservice "github.com/dionis-the-shark/apms-task-tracker/internal/service/auth"
	developerskillservice "github.com/dionis-the-shark/apms-task-tracker/internal/service/developer_skill"
	projectservice "github.com/dionis-the-shark/apms-task-tracker/internal/service/project"
	skillservice "github.com/dionis-the-shark/apms-task-tracker/internal/service/skill"
	taskservice "github.com/dionis-the-shark/apms-task-tracker/internal/service/task"
	taskdependencyservice "github.com/dionis-the-shark/apms-task-tracker/internal/service/task_dependency"
	userservice "github.com/dionis-the-shark/apms-task-tracker/internal/service/user"
)

type Services struct {
	Auth             *authservice.Service
	Users            *userservice.Service
	Projects         *projectservice.Service
	Tasks            *taskservice.Service
	Skills           *skillservice.Service
	DeveloperSkills  *developerskillservice.Service
	TaskDependencies *taskdependencyservice.Service
}

func NewServices(db *sql.DB) *Services {
	authRepo := &repository.AuthRepository{DB: db}
	userRepo := &repository.UserRepository{DB: db}
	projectRepo := &repository.ProjectRepository{DB: db}
	taskRepo := &repository.TaskRepository{DB: db}
	skillRepo := &repository.SkillRepository{DB: db}
	developerSkillRepo := &repository.DeveloperSkillRepository{DB: db}
	taskDependencyRepo := &repository.TaskDependencyRepository{DB: db}

	authService := authservice.New(authRepo)
	userService := userservice.New(userRepo)
	projectService := projectservice.New(projectRepo, userService)
	taskService := taskservice.New(taskRepo)
	skillService := skillservice.New(skillRepo)
	developerSkillService := developerskillservice.New(developerSkillRepo)
	taskDependencyService := taskdependencyservice.New(taskDependencyRepo)

	return &Services{
		Auth:             authService,
		Users:            userService,
		Projects:         projectService,
		Tasks:            taskService,
		Skills:           skillService,
		DeveloperSkills:  developerSkillService,
		TaskDependencies: taskDependencyService,
	}
}
