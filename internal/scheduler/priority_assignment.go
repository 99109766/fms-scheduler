package scheduler

import (
	"github.com/99109766/fms-scheduler/internal/tasks"
	"sort"
)

// DeterminePriorityLevels is a placeholder that, for example,
// might assign priority based on period (rate monotonic) or deadline (EDF).
// For Phase 1, we can do something simple: sort tasks by ascending period (RM).
func DeterminePriorityLevels(taskSet []*tasks.Task) {
	// Sort tasks by ascending period
	sort.Slice(taskSet, func(i, j int) bool {
		return taskSet[i].Period < taskSet[j].Period
	})

	// Assign priority: 1 = highest, 2 = next, ...
	for rank, t := range taskSet {
		// rank starts from 0, so priority is rank+1
		t.UpdatePriority(rank + 1)
	}

	// NOTE: For EDF, you'd sort by absolute deadlines, not by period,
	// or manage a dynamic update in a real-time kernel loop, etc.
}
