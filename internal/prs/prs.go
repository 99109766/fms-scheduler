package prs

import (
	"math"

	"github.com/99109766/fms-scheduler/internal/resources"
	"github.com/99109766/fms-scheduler/internal/tasks"
)

// ComputeResourceCeilings computes and sets the ceiling for each resource.
// The ceiling is defined as the highest priority (i.e. lowest numerical value)
// among the tasks that are assigned to the resource.
func ComputeResourceCeilings(resourceList []*resources.Resource, taskSet []*tasks.Task) {
	// Build a map for quick task lookup by ID.
	taskMap := make(map[int]*tasks.Task)
	for _, t := range taskSet {
		taskMap[t.ID] = t
	}
	for _, r := range resourceList {
		ceiling := math.MaxInt32 // a very high number
		for _, taskID := range r.AssignedTasks {
			if t, ok := taskMap[taskID]; ok {
				if t.Priority < ceiling {
					ceiling = t.Priority
				}
			}
		}
		// If no task uses this resource, assign a default high ceiling.
		if ceiling == math.MaxInt32 {
			ceiling = 9999
		}
		r.Ceiling = ceiling
	}
}

// AssignPreemptionLevels assigns each task a preemption level.
// For a task, the preemption level is defined as the minimum of its base priority
// and the ceilings of all resources it uses.
func AssignPreemptionLevels(taskSet []*tasks.Task, resourceList []*resources.Resource) {
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

// InitPRS is a convenience function that initializes the PRS:
// it computes resource ceilings and then assigns preemption levels.
func InitPRS(taskSet []*tasks.Task, resourceList []*resources.Resource) {
	ComputeResourceCeilings(resourceList, taskSet)
	AssignPreemptionLevels(taskSet, resourceList)
}
