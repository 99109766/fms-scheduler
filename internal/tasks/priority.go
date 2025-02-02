package tasks

import (
	"sort"
)

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
