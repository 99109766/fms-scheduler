package scheduler

import "github.com/99109766/fms-scheduler/internal/tasks"

type Schedule struct {
	TaskID    int     `json:"task_id"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
}

type Job struct {
	Task             *tasks.Task
	JobID            int
	ReleaseTime      float64
	AbsoluteDeadline float64
	RemainingTime    float64
	ExecTime         float64
}

// getActiveCriticalSection returns the active critical section for the job,
// if any. (The one with the shortest duration in case of multiple overlapping CSs.)
func (job *Job) getActiveCriticalSection() *tasks.CriticalSection {
	var bestCS *tasks.CriticalSection
	for _, cs := range job.Task.CriticalSections {
		// Check if job execution is within the CS interval.
		if job.ExecTime >= cs.Start && job.ExecTime < cs.Start+cs.Duration {
			// If the current CS is shorter than the best one, update the best.
			if bestCS == nil || cs.Duration < bestCS.Duration {
				bestCS = cs
			}
		}
	}
	return bestCS
}

// effectivePriority returns a numeric “priority” for the job.
// For jobs not in a critical section we use the absolute deadline (lower is better).
// When inside a critical section the job’s effective priority is its preemption level.
func (job *Job) effectivePriority() float64 {
	if job.getActiveCriticalSection() != nil {
		// When in a critical section, the job’s effective priority is its preemption level.
		return float64(job.Task.PreemptionLevel)
	}
	return job.AbsoluteDeadline
}
