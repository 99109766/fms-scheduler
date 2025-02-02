package tasks

import (
	"math"
	"math/rand"
	"sort"

	"github.com/99109766/fms-scheduler/config"
	"github.com/99109766/fms-scheduler/internal/resources"
)

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
		if len(t.AssignedResIDs) == 0 {
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
			t.CriticalSections = append(t.CriticalSections, &CriticalSection{
				ResourceID: resID,
				Start:      currentTime,
				Duration:   durations[i],
			})
			currentTime += durations[i] + gaps[i+1]
		}
	}
}

// DeterminePriorityLevels assigns static priorities to tasks.
// Here we sort tasks by ascending period (Rate-Monotonic) and assign priorities
// such that lower numbers mean higher priority.
func DeterminePriorityLevels(taskSet []*Task) {
	sort.Slice(taskSet, func(i, j int) bool {
		return taskSet[i].Period < taskSet[j].Period
	})
	for rank, t := range taskSet {
		t.Priority = rank + 1
	}
}

// ComputeResourceCeilings computes and sets the ceiling for each resource.
// The ceiling is defined as the highest priority (i.e. lowest numerical value)
// among the tasks that are assigned to the resource.
func ComputeResourceCeilings(taskSet []*Task, resourceList []*resources.Resource) {
	// Build a map for quick task lookup by ID.
	taskMap := make(map[int]*Task)
	for _, t := range taskSet {
		taskMap[t.ID] = t
	}

	for _, r := range resourceList {
		ceiling := math.MaxInt32
		for _, taskID := range r.AssignedTasks {
			if t, ok := taskMap[taskID]; ok {
				if t.Priority < ceiling {
					ceiling = t.Priority
				}
			}
		}
		r.Ceiling = ceiling
	}
}

// AssignPreemptionLevels assigns each task a preemption level.
// For a task, the preemption level is defined as the minimum of its base priority
// and the ceilings of all resources it uses.
func AssignPreemptionLevels(taskSet []*Task, resourceList []*resources.Resource) {
	// Build a map for quick resource lookup by ID.
	resourceMap := make(map[int]*resources.Resource)
	for _, r := range resourceList {
		resourceMap[r.ID] = r
	}

	for _, t := range taskSet {
		preemptionLevel := t.Priority
		for _, resID := range t.AssignedResIDs {
			if r, ok := resourceMap[resID]; ok {
				if r.Ceiling < preemptionLevel {
					preemptionLevel = r.Ceiling
				}
			}
		}
		t.PreemptionLevel = preemptionLevel
	}
}
