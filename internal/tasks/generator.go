package tasks

import (
	"math"
	"math/rand"

	"github.com/99109766/fms-scheduler/config"
	"github.com/99109766/fms-scheduler/internal/resources"
)

// GenerateTasksUUnifast generates a set of tasks whose sum of utilization = totalUtil.
func GenerateTasksUUnifast(cfg *config.Config) []*Task {
	numTasks, totalUtil := cfg.NumTasks, cfg.TotalUtility

	// Apply UUnifast algorithm
	utilizations := uUniFast(numTasks, totalUtil)

	// Create Task structures
	tasks := make([]*Task, numTasks)
	for i := 0; i < numTasks; i++ {
		period := cfg.PeriodRange[0] + rand.Float64()*(cfg.PeriodRange[1]-cfg.PeriodRange[0])
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

	return tasks
}

// assignRandomCriticalities assigns random criticalities to tasks.
func assignRandomCriticalities(cfg *config.Config, tasks []*Task) {
	// Assign random criticalities (some LC, some HC)
	// and for HC tasks, assign two different WCET values
	for _, t := range tasks {
		if rand.Float64() < cfg.HighRatio {
			t.Criticality = HC
			t.WCET2 = (cfg.WCETRatio[0] + rand.Float64()*(cfg.WCETRatio[1]-cfg.WCETRatio[0])) * t.WCET1
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
			if rand.Float64() < cfg.ResourceUsage {
				t.AssignedResIDs = append(t.AssignedResIDs, r.ID)
				r.AssignedTasks = append(r.AssignedTasks, t.ID)
			}
		}
	}
}

// AssignCriticalSections simulates that each assigned resource has a critical section in the task.
func AssignCriticalSections(cfg *config.Config, tasks []*Task, resources []*resources.Resource) {
	for _, t := range tasks {
		t.CriticalSections = nil
		if t.AssignedResIDs == nil {
			continue
		}

		totalDuration := t.WCET1 * rand.Float64() * cfg.CSFactor
		durations := uUniFast(len(t.AssignedResIDs), totalDuration)
		for i, resID := range t.AssignedResIDs {
			t.CriticalSections = append(t.CriticalSections, CriticalSection{
				ResourceID: resID,
				Duration:   durations[i],
			})
		}
	}
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
