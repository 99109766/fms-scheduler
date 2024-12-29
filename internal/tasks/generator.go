package tasks

import (
	"errors"
	"math"
	"math/rand"

	"github.com/99109766/fms-scheduler/internal/resources"
)

const (
	MinPeriod     = 50  // Period range for tasks
	MaxPeriod     = 200 // Period range for tasks
	WCETRatio     = 0.5 // WCET2 = WCET1 * WCETRatio
	HighRatio     = 0.4 // Probability of a task being high-criticality
	ResourceRatio = 0.5 // Probability of a task using a resource
)

// GenerateTasksUUnifast generates a set of tasks whose sum of utilization = totalUtil.
// numTasks is the count of tasks, totalUtil is the sum of utilization for all tasks.
func GenerateTasksUUnifast(numTasks int, totalUtil float64) ([]*Task, error) {
	if numTasks <= 0 {
		return nil, errors.New("numTasks must be positive")
	}
	if totalUtil <= 0 || totalUtil > float64(numTasks) {
		return nil, errors.New("totalUtil is out of valid range")
	}

	// Apply UUnifast algorithm
	utilizations := uUniFast(numTasks, totalUtil)

	// Create Task structures
	tasks := make([]*Task, numTasks)
	for i := 0; i < numTasks; i++ {
		period := MinPeriod + rand.Float64()*(MaxPeriod-MinPeriod)
		wcet := utilizations[i] * period

		tasks[i] = &Task{
			ID:       i,
			WCET1:    wcet,
			Period:   period,
			Deadline: period, // By default, deadline = period
		}
	}

	// Assign random criticalities (some LC, some HC)
	// and for HC tasks, assign two different WCET values
	for _, t := range tasks {
		if rand.Float64() < HighRatio {
			t.Criticality = HC
			t.WCET2 = rand.Float64() * WCETRatio * t.WCET1
		} else {
			t.Criticality = LC
			t.WCET2 = 0
		}
	}

	return tasks, nil
}

// AssignResourcesToTasks randomly assigns resources to tasks.
func AssignResourcesToTasks(tasks []*Task, resources []*resources.Resource) {
	for _, r := range resources {
		r.AssignedTasks = nil
	}

	for _, t := range tasks {
		t.AssignedResIDs = nil
		for _, r := range resources {
			if rand.Float64() < ResourceRatio {
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
	// We'll keep it no-op here or minimal for phase 1.
	// TODO: complete this functio
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
