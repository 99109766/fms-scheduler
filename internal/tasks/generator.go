package tasks

import (
	"errors"
	"math"
	"math/rand"

	"github.com/99109766/fms-scheduler/config"
	"github.com/99109766/fms-scheduler/internal/resources"
)

// GenerateTasksUUnifast generates a set of tasks whose sum of utilization = totalUtil.
func GenerateTasksUUnifast(cfg *config.Config) ([]*Task, error) {
	numTasks, totalUtil := cfg.NumTasks, cfg.TotalUtil
	if numTasks <= 0 {
		return nil, errors.New("number of tasks must be greater than 0")
	}
	if totalUtil <= 0 || totalUtil > 1 {
		return nil, errors.New("total utilization must be in the range (0, 1]")
	}

	// Apply UUnifast algorithm
	utilizations := uUniFast(numTasks, totalUtil)

	// Create Task structures
	tasks := make([]*Task, numTasks)
	for i := 0; i < numTasks; i++ {
		period := cfg.MinPeriod + rand.Float64()*(cfg.MaxPeriod-cfg.MinPeriod)
		wcet := utilizations[i] * period

		tasks[i] = &Task{
			ID:       i,
			WCET1:    wcet,
			Period:   period,
			Deadline: period, // By default, deadline = period
		}
	}

	// Assign random criticalities
	assignRandomCriticalities(cfg, tasks)

	return tasks, nil
}

// assignRandomCriticalities assigns random criticalities to tasks.
func assignRandomCriticalities(cfg *config.Config, tasks []*Task) {
	// Assign random criticalities (some LC, some HC)
	// and for HC tasks, assign two different WCET values
	for _, t := range tasks {
		if rand.Float64() < cfg.HighRatio {
			t.Criticality = HC
			t.WCET2 = rand.Float64() * cfg.WCETRatio * t.WCET1
		} else {
			t.Criticality = LC
			t.WCET2 = 0
		}
	}
}

// AssignResourcesToTasks randomly assigns resources to tasks.
func AssignResourcesToTasks(cfg *config.Config, tasks []*Task, resources []*resources.Resource) {
	for _, r := range resources {
		r.AssignedTasks = nil
	}

	for _, t := range tasks {
		t.AssignedResIDs = nil
		for _, r := range resources {
			if rand.Float64() < cfg.ResourceRatio {
				t.AssignedResIDs = append(t.AssignedResIDs, r.ID)
				r.AssignedTasks = append(r.AssignedTasks, t.ID)
			}
		}
	}
}

// AssignCriticalSections simulates that each assigned resource has a critical section in the task.
func AssignCriticalSections(tasks []*Task, resourceList []*resources.Resource) {
	// Here, you could store additional info about critical sections in the tasks
	// (like how long each critical section is). For demonstration, we just do a print or log.
	// You could also integrate a maximum concurrency parameter, etc.
}

// uUniFast is the internal function implementing the UUnifast algorithm.
func uUniFast(n int, U float64) []float64 {
	var sumU float64 = U
	var utilis = make([]float64, n)

	for i := 1; i < n; i++ {
		next := sumU * math.Pow(rand.Float64(), 1.0/float64(n-i))
		utilis[i-1] = sumU - next
		sumU = next
	}
	utilis[n-1] = sumU
	return utilis
}
