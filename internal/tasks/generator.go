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

	assignRandomCriticality(cfg, tasks)

	return tasks
}

// assignRandomCriticality assigns random criticality to tasks.
func assignRandomCriticality(cfg *config.Config, tasks []*Task) {
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
// The critical sections are assigned start times and durations so that they do not partially overlap.
// They are “stacked” sequentially (which is acceptable since non-overlap implies they are not half‐overlapping).
func AssignCriticalSections(cfg *config.Config, tasks []*Task, resources []*resources.Resource) {
	for _, t := range tasks {
		t.CriticalSections = nil
		if t.AssignedResIDs == nil || len(t.AssignedResIDs) == 0 {
			continue
		}

		// Total critical section duration (a fraction of WCET1)
		totalDuration := t.WCET1 * rand.Float64() * cfg.CSFactor

		// Split totalDuration among the assigned resources using uUniFast
		durations := uUniFast(len(t.AssignedResIDs), totalDuration)

		// Compute available free time in the task (WCET1 minus total CS duration)
		freeTime := t.WCET1 - totalDuration
		if freeTime < 0 {
			freeTime = 0
		}
		// Distribute free time as gaps before, between, and after critical sections.
		gaps := uUniFast(len(t.AssignedResIDs)+1, freeTime)

		// Place critical sections sequentially.
		currentTime := gaps[0]
		for i, resID := range t.AssignedResIDs {
			cs := CriticalSection{
				ResourceID: resID,
				Start:      currentTime,
				Duration:   durations[i],
			}
			t.CriticalSections = append(t.CriticalSections, cs)
			currentTime += durations[i] + gaps[i+1]
		}
	}
}

// uUniFast is the internal function implementing the UUniFast algorithm.
func uUniFast(n int, U float64) []float64 {
	sumU := U
	utils := make([]float64, n)
	for i := 1; i < n; i++ {
		next := sumU * (math.Pow(rand.Float64(), 1.0/float64(n-i)))
		utils[i-1] = sumU - next
		sumU = next
	}
	utils[n-1] = sumU
	return utils
}
