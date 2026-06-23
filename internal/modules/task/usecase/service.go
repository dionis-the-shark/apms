package usecase

import (
	"context"
	"errors"
	"sort"
	"strings"
	"time"

	developerskillports "github.com/dionis-the-shark/apms-task-tracker/internal/modules/developer_skill/ports"
	taskmodule "github.com/dionis-the-shark/apms-task-tracker/internal/modules/task"
	"github.com/dionis-the-shark/apms-task-tracker/internal/modules/task/ports"
	userports "github.com/dionis-the-shark/apms-task-tracker/internal/modules/user/ports"
	"github.com/google/uuid"
)

type Service struct {
	repo               ports.Repository
	userRepo           userports.Repository
	developerSkillRepo developerskillports.Repository
	now                func() time.Time
	newUUID            func() uuid.UUID
}

type API interface {
	Create(ctx context.Context, input CreateInput) (taskmodule.Task, error)
	GetByID(ctx context.Context, taskID uuid.UUID) (taskmodule.Task, error)
	GetByProject(ctx context.Context, projectID uuid.UUID) ([]taskmodule.Task, error)
	Update(ctx context.Context, taskID uuid.UUID, input UpdateInput) (taskmodule.Task, error)
	Delete(ctx context.Context, taskID uuid.UUID) error
	DistributeFreeTasks(ctx context.Context, projectID uuid.UUID) (DistributionResult, error)
}

type CreateInput struct {
	ProjectID       uuid.UUID  `json:"project_id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description"`
	Status          string     `json:"status"`
	Priority        int        `json:"priority"`
	EstimatedHours  int        `json:"estimated_hours"`
	RequiredSkillID *uuid.UUID `json:"required_skill_id"`
	ExecutorID      *uuid.UUID `json:"executor_id"`
}

type UpdateInput struct {
	ProjectID       uuid.UUID  `json:"project_id"`
	Title           string     `json:"title"`
	Description     *string    `json:"description"`
	Status          string     `json:"status"`
	Priority        int        `json:"priority"`
	EstimatedHours  int        `json:"estimated_hours"`
	RequiredSkillID *uuid.UUID `json:"required_skill_id"`
	ExecutorID      *uuid.UUID `json:"executor_id"`
}

type DistributionResult struct {
	Assigned int               `json:"assigned"`
	Skipped  int               `json:"skipped"`
	Tasks    []taskmodule.Task `json:"tasks"`
}

func New(repo ports.Repository, userRepo userports.Repository, developerSkillRepo developerskillports.Repository) *Service {
	return &Service{
		repo:               repo,
		userRepo:           userRepo,
		developerSkillRepo: developerSkillRepo,
		now:                func() time.Time { return time.Now().UTC() },
		newUUID:            uuid.New,
	}
}

func (s *Service) Create(ctx context.Context, input CreateInput) (taskmodule.Task, error) {
	status := input.Status
	if status == "" {
		status = "backlog"
	}
	priority := input.Priority
	if priority == 0 {
		priority = 2
	}

	t := taskmodule.Task{
		TaskID:          s.newUUID(),
		ProjectID:       input.ProjectID,
		Title:           input.Title,
		Description:     input.Description,
		Status:          status,
		Priority:        priority,
		EstimatedHours:  input.EstimatedHours,
		RequiredSkillID: input.RequiredSkillID,
		ExecutorID:      input.ExecutorID,
		CreatedAt:       s.now(),
	}
	if err := s.repo.CreateTask(ctx, &t); err != nil {
		return taskmodule.Task{}, err
	}
	return t, nil
}

func (s *Service) GetByID(ctx context.Context, taskID uuid.UUID) (taskmodule.Task, error) {
	return s.repo.GetTaskByID(ctx, taskID)
}

func (s *Service) GetByProject(ctx context.Context, projectID uuid.UUID) ([]taskmodule.Task, error) {
	return s.repo.GetTasksByProject(ctx, projectID)
}

func (s *Service) Update(ctx context.Context, taskID uuid.UUID, input UpdateInput) (taskmodule.Task, error) {
	status := input.Status
	if status == "" {
		status = "backlog"
	}
	priority := input.Priority
	if priority == 0 {
		priority = 2
	}

	t := taskmodule.Task{
		TaskID:          taskID,
		ProjectID:       input.ProjectID,
		Title:           input.Title,
		Description:     input.Description,
		Status:          status,
		Priority:        priority,
		EstimatedHours:  input.EstimatedHours,
		RequiredSkillID: input.RequiredSkillID,
		ExecutorID:      input.ExecutorID,
	}
	if err := s.repo.UpdateTask(ctx, t); err != nil {
		return taskmodule.Task{}, err
	}
	return t, nil
}

func (s *Service) Delete(ctx context.Context, taskID uuid.UUID) error {
	return s.repo.DeleteTask(ctx, taskID)
}

type workload struct {
	hours    int
	total    int
	priority map[int]int
}

func (w *workload) addTask(priority int, hours int) {
	if w.priority == nil {
		w.priority = make(map[int]int)
	}
	w.priority[priority]++
	w.hours += hours
	w.total++
}

func (w *workload) score() float64 {
	// Weighted score: hours + weighted task count by priority.
	weights := map[int]float64{
		1: 1,
		2: 2,
		3: 3,
		4: 5,
	}
	score := float64(w.hours)
	for priority, count := range w.priority {
		score += weights[priority] * float64(count)
	}
	return score
}

func normalizeStatus(status string) string {
	// Normalize status values from various inputs to canonical buckets.
	raw := strings.ToLower(strings.TrimSpace(status))
	raw = strings.ReplaceAll(raw, " ", "_")
	switch {
	case strings.Contains(raw, "progress"):
		return "in_progress"
	case strings.Contains(raw, "test"):
		return "testing"
	case strings.Contains(raw, "done"):
		return "done"
	default:
		return "backlog"
	}
}

// compareWorkload compares workloads of two developers. Returns:
//  1. -1 if dev1 workload is lower than dev1 workload,
//  2. 1 if dev1 workload is higher than dev1 workload,
//  3. 0 if equal.
func compareWorkload(dev1, dev2 *workload) int {
	scoreA := dev1.score()
	scoreB := dev2.score()
	if scoreA < scoreB {
		return -1
	}
	if scoreA > scoreB {
		return 1
	}
	priorityOrder := []int{4, 3, 2, 1}
	for _, p := range priorityOrder {
		if dev1.priority[p] < dev2.priority[p] {
			return -1
		}
		if dev1.priority[p] > dev2.priority[p] {
			return 1
		}
	}
	if dev1.hours < dev2.hours {
		return -1
	}
	if dev1.hours > dev2.hours {
		return 1
	}
	if dev1.total < dev2.total {
		return -1
	}
	if dev1.total > dev2.total {
		return 1
	}
	return 0
}

func (s *Service) DistributeFreeTasks(ctx context.Context, projectID uuid.UUID) (DistributionResult, error) {
	if s.userRepo == nil || s.developerSkillRepo == nil {
		return DistributionResult{}, errors.New("distribution dependencies are not configured")
	}

	// Load project tasks, users, and developer skills.
	tasks, err := s.repo.GetTasksByProject(ctx, projectID)
	if err != nil {
		return DistributionResult{}, err
	}
	users, err := s.userRepo.GetUsers(ctx)
	if err != nil {
		return DistributionResult{}, err
	}
	developerSkills, err := s.developerSkillRepo.GetDeveloperSkills(ctx)
	if err != nil {
		return DistributionResult{}, err
	}

	// Build user set and skill lookup maps.
	userSet := make(map[uuid.UUID]struct{})
	for _, user := range users {
		userSet[user.UserID] = struct{}{}
	}

	developersBySkill := make(map[uuid.UUID][]uuid.UUID)
	developerSkillSet := make(map[uuid.UUID]map[uuid.UUID]struct{})
	for _, ds := range developerSkills {
		if _, ok := userSet[ds.DeveloperID]; !ok {
			continue
		}
		developersBySkill[ds.SkillID] = append(developersBySkill[ds.SkillID], ds.DeveloperID)
		if developerSkillSet[ds.DeveloperID] == nil {
			developerSkillSet[ds.DeveloperID] = make(map[uuid.UUID]struct{})
		}
		developerSkillSet[ds.DeveloperID][ds.SkillID] = struct{}{}
	}

	// Precompute candidate developer lists.
	allUsers := make([]uuid.UUID, 0, len(userSet))
	for id := range userSet {
		allUsers = append(allUsers, id)
	}

	developers := make([]uuid.UUID, 0, len(developerSkillSet))
	for devID := range developerSkillSet {
		developers = append(developers, devID)
	}

	// Calculate current workload from all non-done assigned tasks.
	workloads := make(map[uuid.UUID]*workload)
	for _, devID := range allUsers {
		workloads[devID] = &workload{priority: make(map[int]int)}
	}

	for _, task := range tasks {
		if task.ExecutorID == nil {
			continue
		}
		if normalizeStatus(task.Status) == "done" {
			continue
		}
		stats, ok := workloads[*task.ExecutorID]
		if !ok {
			stats = &workload{priority: make(map[int]int)}
			workloads[*task.ExecutorID] = stats
		}
		stats.addTask(task.Priority, task.EstimatedHours)
	}

	// Gather free backlog tasks only.
	freeTasks := make([]taskmodule.Task, 0)
	for _, task := range tasks {
		if normalizeStatus(task.Status) != "backlog" {
			continue
		}
		if task.ExecutorID != nil {
			continue
		}
		freeTasks = append(freeTasks, task)
	}

	// Distribute higher-priority and larger tasks first.
	sort.SliceStable(freeTasks, func(i, j int) bool {
		if freeTasks[i].Priority != freeTasks[j].Priority {
			return freeTasks[i].Priority > freeTasks[j].Priority
		}
		return freeTasks[i].EstimatedHours > freeTasks[j].EstimatedHours
	})

	result := DistributionResult{}

	for _, task := range freeTasks {
		// Determine eligible developers: skill match if required, otherwise any developer/user.
		var eligible []uuid.UUID
		if task.RequiredSkillID != nil {
			eligible = append(eligible, developersBySkill[*task.RequiredSkillID]...)
		} else if len(developers) > 0 {
			eligible = append(eligible, developers...)
		} else {
			eligible = append(eligible, allUsers...)
		}

		if len(eligible) == 0 {
			continue
		}

		// Pick the least-loaded eligible developer.
		sort.SliceStable(eligible, func(i, j int) bool {
			a := workloads[eligible[i]]
			b := workloads[eligible[j]]
			if a == nil {
				a = &workload{priority: make(map[int]int)}
			}
			if b == nil {
				b = &workload{priority: make(map[int]int)}
			}
			compare := compareWorkload(a, b)
			if compare != 0 {
				return compare < 0
			}
			return eligible[i].String() < eligible[j].String()
		})

		chosen := eligible[0]
		executorID := chosen
		task.ExecutorID = &executorID
		// Persist assignment and update in-memory workload.
		if err := s.repo.UpdateTask(ctx, task); err != nil {
			return result, err
		}
		stats := workloads[chosen]
		if stats == nil {
			stats = &workload{priority: make(map[int]int)}
			workloads[chosen] = stats
		}
		stats.addTask(task.Priority, task.EstimatedHours)
		result.Assigned++
		result.Tasks = append(result.Tasks, task)
	}

	result.Skipped = len(freeTasks) - result.Assigned
	return result, nil
}
