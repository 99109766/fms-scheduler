package tasks

import (
	"fmt"
)

// CriticalityLevel represents the criticality of a task (LC or HC).
type CriticalityLevel int

const (
	LC CriticalityLevel = iota
	HC
)

// Task models a two-level mixed-criticality periodic task.
type Task struct {
	ID          int
	Criticality CriticalityLevel
	Period      float64
	Deadline    float64
	// For Low-Criticality tasks (LC), only WCET1 is used.
	// For High-Criticality tasks (HC), we have WCET1 and WCET2.
	WCET1          float64
	WCET2          float64 // only used if HC
	Utilization    float64 // For convenience: WCET/Period
	AssignedResIDs []int   // Resource IDs this task will use
	Priority       int     // Priority (lower number = higher priority, example)
}

// String implements fmt.Stringer for pretty print
func (t Task) String() string {
	if t.Criticality == LC {
		return fmt.Sprintf(
			"[Task %d | LC | Period=%.2f | Deadline=%.2f | WCET=%.2f | Util=%.2f | Priority=%d | Res=%v]",
			t.ID, t.Period, t.Deadline, t.WCET1, t.Utilization, t.Priority, t.AssignedResIDs)
	}
	return fmt.Sprintf(
		"[Task %d | HC | Period=%.2f | Deadline=%.2f | WCET1=%.2f | WCET2=%.2f | Util=%.2f | Priority=%d | Res=%v]",
		t.ID, t.Period, t.Deadline, t.WCET1, t.WCET2, t.Utilization, t.Priority, t.AssignedResIDs)
}

// ComputeUtilization updates the task’s utilization field based on its WCET & Period
func (t *Task) ComputeUtilization() {
	// For LC tasks or the normal operation mode, we consider WCET1
	t.Utilization = t.WCET1 / t.Period
	// If you need to consider worst-case (overrun) scenarios for HC tasks,
	// you could store that in a separate field or keep it for the scheduling phase.
}

// UpdatePriority is a helper for priority assignment
func (t *Task) UpdatePriority(priority int) {
	t.Priority = priority
}

// ComputeDeadline can be used if you want to align tasks’ deadlines with periods
// or do something custom. Here, assume deadline = period for simplicity.
func (t *Task) ComputeDeadline() {
	t.Deadline = t.Period
}
