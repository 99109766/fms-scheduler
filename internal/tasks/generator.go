package tasks

import (
	"errors"
	"github.com/99109766/fms-scheduler/internal/resources"
	"math"
	"math/rand"
	"time"
)

// GenerateTasksUUnifast generates a set of tasks whose sum of utilization = totalUtil.
// numTasks is the count of tasks, totalUtil is the sum of utilization for all tasks.
func GenerateTasksUUnifast(numTasks int, totalUtil float64) ([]*Task, error) {
	if numTasks <= 0 {
		return nil, errors.New("numTasks must be > 0")
	}
	if totalUtil <= 0 || totalUtil > float64(numTasks) {
		return nil, errors.New("totalUtil is out of valid range")
	}

	rand.Seed(time.Now().UnixNano())

	// Apply UUnifast algorithm
	utilizations := uUniFast(numTasks, totalUtil)

	// Create Task structures
	tasks := make([]*Task, numTasks)
	for i := 0; i < numTasks; i++ {
		// Example period range: [50, 200]
		period := float64(rand.Intn(150) + 50)
		wcet := utilizations[i] * period

		tasks[i] = &Task{
			ID:     i,
			Period: period,
			WCET1:  wcet,
		}
		// By default, let's align the deadline with period:
		tasks[i].ComputeDeadline()
		tasks[i].ComputeUtilization()
	}

	return tasks, nil
}

// AssignMixedCriticality randomly selects some tasks to be High-Criticality and sets WCET2
func AssignMixedCriticality(tasks []*Task) {
	for _, t := range tasks {
		if rand.Float64() < 0.4 {
			t.Criticality = HC
			// Suppose WCET2 is 1.5 to 2 times bigger than WCET1
			t.WCET2 = t.WCET1 * (1.5 + 0.5*rand.Float64())
		} else {
			t.Criticality = LC
			// For LC tasks, WCET2 is not used, keep it zero or consistent
			t.WCET2 = 0
		}
	}
}

// AssignResourcesToTasks randomly assigns resources to tasks
func AssignResourcesToTasks(tasks []*Task, resources []resources.Resource) {
	// For demonstration, each task might use a subset of the resource IDs
	for _, t := range tasks {
		t.AssignedResIDs = nil
		// Randomly decide if we use each resource
		for _, r := range resources {
			if rand.Float64() < 0.5 {
				t.AssignedResIDs = append(t.AssignedResIDs, r.ID)
			}
		}
	}
}

// AssignCriticalSections simulates that each assigned resource has a critical section in the task
// This is a placeholder: in a real system, you'd define code sections, lengths, etc.
func AssignCriticalSections(tasks []*Task, resourceList []resources.Resource) {
	// Here, you could store additional info about critical sections in the tasks
	// (like how long each critical section is). For demonstration, we just do a print or log.
	// You could also integrate a maximum concurrency parameter, etc.
	// We'll keep it no-op here or minimal for phase 1.
	// TODO: complete this function
}

// uUniFast is the internal function implementing the UUnifast algorithm.
func uUniFast(n int, U float64) []float64 {
	var sumU float64 = U
	var utilis = make([]float64, n)

	for i := 1; i < n; i++ {
		// draw a random number
		next := sumU * math.Pow(rand.Float64(), 1.0/float64(n-i))
		utilis[i-1] = sumU - next
		sumU = next
	}
	utilis[n-1] = sumU
	return utilis
}
