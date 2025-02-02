package scheduler

import (
	"fmt"
	"sort"

	"github.com/99109766/fms-scheduler/internal/tasks"
)

// Mode indicates the current system mode.
type Mode int

const (
	Normal Mode = iota
	Overrun
)

// Job models a released instance of a task.
type Job struct {
	Task             *tasks.Task // pointer to the originating task
	JobID            int         // unique job identifier
	ReleaseTime      float64     // when the job is released
	AbsoluteDeadline float64     // release time + task deadline
	RemainingTime    float64     // how much execution remains for this job
	ExecTime         float64     // cumulative execution time so far
}

// jobCounter is used to assign unique IDs to jobs.
var jobCounter int = 0

// inCriticalSection returns true if the job’s current execution time falls
// into one of the task’s critical section intervals.
func inCriticalSection(job *Job) bool {
	for _, cs := range job.Task.CriticalSections {
		// Check if job execution is within the CS interval.
		if job.ExecTime >= cs.Start && job.ExecTime < cs.Start+cs.Duration {
			return true
		}
	}
	return false
}

// getActiveCriticalSection returns the active critical section for the job,
// if any. (Since our assignment places CS intervals sequentially, only one
// will be active at any given time.)
func getActiveCriticalSection(job *Job) *tasks.CriticalSection {
	for i := range job.Task.CriticalSections {
		cs := &job.Task.CriticalSections[i]
		if job.ExecTime >= cs.Start && job.ExecTime < cs.Start+cs.Duration {
			return cs
		}
	}
	return nil
}

// effectivePriority returns a numeric “priority” for the job.
// For jobs not in a critical section we use the absolute deadline (lower is better).
// When inside a critical section the job’s effective priority is its preemption level.
func effectivePriority(job *Job) float64 {
	if inCriticalSection(job) {
		// When in a critical section, the job’s effective priority is its preemption level.
		return float64(job.Task.PreemptionLevel)
	}
	return job.AbsoluteDeadline
}

// RunScheduler simulates an EDF-ER scheduler for a mixed-criticality system.
// It releases jobs from the task set and simulates execution for simTime seconds.
func RunScheduler(taskSet []*tasks.Task, simTime float64) {
	mode := Normal
	currentTime := 0.0
	// Use a small simulation step (dt). (A real system would be event-driven.)
	dt := 0.001

	// Map each task to its next release time.
	nextRelease := make(map[int]float64)
	for _, t := range taskSet {
		nextRelease[t.ID] = 0.0
	}

	// readyQueue holds released jobs waiting to run.
	readyQueue := make([]*Job, 0)
	var runningJob *Job = nil

	// scheduleLog records scheduling events.
	scheduleLog := []string{}

	// Track whether the running job was in a critical section in the previous dt step.
	runningJobInCS := false

	// Main simulation loop.
	for currentTime < simTime {

		// 1. Release new jobs if their release time has arrived.
		for _, t := range taskSet {
			if currentTime >= nextRelease[t.ID] {
				// In Overrun mode, only release HC tasks.
				if mode == Normal || (mode == Overrun && t.Criticality == tasks.HC) {
					newJob := &Job{
						Task:             t,
						ReleaseTime:      nextRelease[t.ID],
						AbsoluteDeadline: nextRelease[t.ID] + t.Deadline,
						ExecTime:         0,
					}
					// For HC tasks, choose WCET1 in Normal mode and WCET2 in Overrun mode.
					if t.Criticality == tasks.HC {
						if mode == Normal {
							newJob.RemainingTime = t.WCET1
						} else {
							newJob.RemainingTime = t.WCET2
						}
					} else {
						newJob.RemainingTime = t.WCET1
					}
					jobCounter++
					newJob.JobID = jobCounter
					readyQueue = append(readyQueue, newJob)
					scheduleLog = append(scheduleLog,
						fmt.Sprintf("Time %.3f: Released Job %d (Task %d, Deadline=%.3f, WCET=%.3f) [Mode: %v]",
							currentTime, newJob.JobID, t.ID, newJob.AbsoluteDeadline, newJob.RemainingTime, mode))
				}
				// Schedule the next release for the task.
				nextRelease[t.ID] += t.Period
			}
		}

		// 2. In Overrun mode, drop any LC jobs from the ready queue.
		if mode == Overrun {
			newReady := []*Job{}
			for _, job := range readyQueue {
				if job.Task.Criticality == tasks.HC {
					newReady = append(newReady, job)
				} else {
					scheduleLog = append(scheduleLog,
						fmt.Sprintf("Time %.3f: Dropped LC Job %d (Task %d) due to Overrun mode",
							currentTime, job.JobID, job.Task.ID))
				}
			}
			readyQueue = newReady
		}

		// 3. Scheduler decision: select a job to run.
		if runningJob == nil && len(readyQueue) > 0 {
			// Pick the job with the smallest effective priority.
			sort.Slice(readyQueue, func(i, j int) bool {
				return effectivePriority(readyQueue[i]) < effectivePriority(readyQueue[j])
			})
			runningJob = readyQueue[0]
			readyQueue = readyQueue[1:]
			scheduleLog = append(scheduleLog,
				fmt.Sprintf("Time %.3f: Starting Job %d (Task %d) with Deadline=%.3f, EffectivePriority=%.3f",
					currentTime, runningJob.JobID, runningJob.Task.ID, runningJob.AbsoluteDeadline, effectivePriority(runningJob)))
			runningJobInCS = inCriticalSection(runningJob)
			if runningJobInCS {
				if cs := getActiveCriticalSection(runningJob); cs != nil {
					scheduleLog = append(scheduleLog,
						fmt.Sprintf("Time %.3f: Job %d ENTERS critical section on Resource %d (CS: Start=%.3f, Duration=%.3f)",
							currentTime, runningJob.JobID, cs.ResourceID, cs.Start, cs.Duration))
				}
			}
		} else if runningJob != nil && len(readyQueue) > 0 {
			// Check if a waiting job has a better (lower) effective priority.
			sort.Slice(readyQueue, func(i, j int) bool {
				return effectivePriority(readyQueue[i]) < effectivePriority(readyQueue[j])
			})
			candidate := readyQueue[0]
			// Preempt if the candidate’s effective priority is better.
			if effectivePriority(candidate) < effectivePriority(runningJob) {
				// If runningJob is in a critical section, allow preemption only if candidate beats its preemption level.
				if inCriticalSection(runningJob) {
					if effectivePriority(candidate) < float64(runningJob.Task.PreemptionLevel) {
						scheduleLog = append(scheduleLog,
							fmt.Sprintf("Time %.3f: Preempting Job %d (Task %d, EffPri=%.3f, in CS) with Job %d (Task %d, EffPri=%.3f)",
								currentTime, runningJob.JobID, runningJob.Task.ID, effectivePriority(runningJob),
								candidate.JobID, candidate.Task.ID, effectivePriority(candidate)))
						readyQueue = append(readyQueue, runningJob)
						runningJob = candidate
						readyQueue = readyQueue[1:]
						runningJobInCS = inCriticalSection(runningJob)
						if runningJobInCS {
							if cs := getActiveCriticalSection(runningJob); cs != nil {
								scheduleLog = append(scheduleLog,
									fmt.Sprintf("Time %.3f: Job %d ENTERS critical section on Resource %d (CS: Start=%.3f, Duration=%.3f)",
										currentTime, runningJob.JobID, cs.ResourceID, cs.Start, cs.Duration))
							}
						}
					}
				} else {
					scheduleLog = append(scheduleLog,
						fmt.Sprintf("Time %.3f: Preempting Job %d (Task %d, EffPri=%.3f) with Job %d (Task %d, EffPri=%.3f)",
							currentTime, runningJob.JobID, runningJob.Task.ID, effectivePriority(runningJob),
							candidate.JobID, candidate.Task.ID, effectivePriority(candidate)))
					readyQueue = append(readyQueue, runningJob)
					runningJob = candidate
					readyQueue = readyQueue[1:]
					runningJobInCS = inCriticalSection(runningJob)
					if runningJobInCS {
						if cs := getActiveCriticalSection(runningJob); cs != nil {
							scheduleLog = append(scheduleLog,
								fmt.Sprintf("Time %.3f: Job %d ENTERS critical section on Resource %d (CS: Start=%.3f, Duration=%.3f)",
									currentTime, runningJob.JobID, cs.ResourceID, cs.Start, cs.Duration))
						}
					}
				}
			}
		}

		// 3.5. Check and log critical section entry/exit transitions.
		if runningJob != nil {
			currentInCS := inCriticalSection(runningJob)
			if !runningJobInCS && currentInCS {
				if cs := getActiveCriticalSection(runningJob); cs != nil {
					scheduleLog = append(scheduleLog,
						fmt.Sprintf("Time %.3f: Job %d (Task %d) ENTERS critical section on Resource %d (CS: Start=%.3f, Duration=%.3f)",
							currentTime, runningJob.JobID, runningJob.Task.ID, cs.ResourceID, cs.Start, cs.Duration))
				}
			} else if runningJobInCS && !currentInCS {
				scheduleLog = append(scheduleLog,
					fmt.Sprintf("Time %.3f: Job %d (Task %d) EXITS critical section", currentTime, runningJob.JobID, runningJob.Task.ID))
			}
			runningJobInCS = currentInCS
		}

		// 4. Execute the running job for one time step.
		if runningJob != nil {
			runningJob.ExecTime += dt
			runningJob.RemainingTime -= dt

			// Check if an HC job overruns its normal (WCET1) execution in Normal mode.
			if runningJob.Task.Criticality == tasks.HC && mode == Normal && runningJob.ExecTime > runningJob.Task.WCET1 {
				mode = Overrun
				scheduleLog = append(scheduleLog,
					fmt.Sprintf("Time %.3f: Mode switch to OVERRUN triggered by Job %d (Task %d) [ExecTime=%.3f, WCET1=%.3f]",
						currentTime, runningJob.JobID, runningJob.Task.ID, runningJob.ExecTime, runningJob.Task.WCET1))
				// Extend the remaining time to allow execution up to WCET2.
				extra := runningJob.Task.WCET2 - runningJob.Task.WCET1
				runningJob.RemainingTime += extra

				// Drop pending LC jobs.
				newReady := []*Job{}
				for _, job := range readyQueue {
					if job.Task.Criticality == tasks.HC {
						newReady = append(newReady, job)
					} else {
						scheduleLog = append(scheduleLog,
							fmt.Sprintf("Time %.3f: Dropped LC Job %d (Task %d) due to mode switch",
								currentTime, job.JobID, job.Task.ID))
					}
				}
				readyQueue = newReady
			}

			// Job completion.
			if runningJob.RemainingTime <= 0 {
				scheduleLog = append(scheduleLog,
					fmt.Sprintf("Time %.3f: COMPLETED Job %d (Task %d) [FinishTime=%.3f, Total ExecTime=%.3f]",
						currentTime, runningJob.JobID, runningJob.Task.ID, currentTime, runningJob.ExecTime))
				runningJob = nil
			}
		}

		currentTime += dt
	}

	// Print the full scheduling log.
	fmt.Println("=== Scheduler Log ===")
	for _, entry := range scheduleLog {
		fmt.Println(entry)
	}
}
