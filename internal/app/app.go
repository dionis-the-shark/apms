package app

import (
	"database/sql"

	authrepo "github.com/dionis-the-shark/apms-task-tracker/internal/modules/auth/repository"
	authusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/auth/usecase"
	developerskillrepo "github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill/repository"
	developerskillusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill/usecase"
	projectrepo "github.com/dionis-the-shark/apms-task-tracker/internal/modules/project/repository"
	projectusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/project/usecase"
	rolerepo "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role/repository"
	roleusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/role/usecase"
	skillrepo "github.com/dionis-the-shark/apms-task-tracker/internal/modules/skill/repository"
	skillusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/skill/usecase"
	taskrepo "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task/repository"
	taskusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task/usecase"
	taskdependencyrepo "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task_dependency/repository"
	taskdependencyusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task_dependency/usecase"
	userrepo "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user/repository"
	userusecase "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user/usecase"
)

type Services struct {
	Auth             authusecase.API
	Users            userusecase.API
	Projects         projectusecase.API
	Tasks            taskusecase.API
	Skills           skillusecase.API
	DeveloperSkills  developerskillusecase.API
	TaskDependencies taskdependencyusecase.API
	Roles            roleusecase.API
}

func NewServices(db *sql.DB) *Services {
	authRepo := &authrepo.Repository{DB: db}
	userRepo := &userrepo.Repository{DB: db}
	projectRepo := &projectrepo.Repository{DB: db}
	taskRepo := &taskrepo.Repository{DB: db}
	skillRepo := &skillrepo.Repository{DB: db}
	developerSkillRepo := &developerskillrepo.Repository{DB: db}
	taskDependencyRepo := &taskdependencyrepo.Repository{DB: db}
	roleRepo := &rolerepo.Repository{DB: db}

	authService := authusecase.New(authRepo)
	userService := userusecase.New(userRepo)
	projectService := projectusecase.New(projectRepo)
	taskService := taskusecase.New(taskRepo)
	skillService := skillusecase.New(skillRepo)
	developerSkillService := developerskillusecase.New(developerSkillRepo)
	taskDependencyService := taskdependencyusecase.New(taskDependencyRepo)
	roleService := roleusecase.New(roleRepo)

	return &Services{
		Auth:             authService,
		Users:            userService,
		Projects:         projectService,
		Tasks:            taskService,
		Skills:           skillService,
		DeveloperSkills:  developerSkillService,
		TaskDependencies: taskDependencyService,
		Roles:            roleService,
	}
}
