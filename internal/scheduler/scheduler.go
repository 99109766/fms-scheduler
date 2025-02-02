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

// dropLCJobs removes all low-criticality jobs from the ready queue.
// This is used when the system switches to Overrun mode.
func dropLCJobs(queue []*Job, currentTime float64) []*Job {
	newQueue := []*Job{}
	for _, job := range queue {
		if job.Task.Criticality == tasks.HC {
			newQueue = append(newQueue, job)
		} else {
			fmt.Printf("Time %.3f: Dropped LC Job %d (Task %d) due to mode switch\n",
				currentTime, job.JobID, job.Task.ID)
		}
	}
	return newQueue
}

// extendRemainingTime extends the remaining time of all HC jobs in the ready queue.
// This is used when the system switches to Overrun mode.
func extendRemainingTime(jobs []*Job) {
	for _, job := range jobs {
		job.RemainingTime += job.Task.WCET2
	}
}

// RunScheduler simulates an ER-EDF scheduler for a mixed-criticality system.
// It releases jobs from the task set and simulates execution for simTime seconds.
func RunScheduler(taskSet []*tasks.Task, simulateTime float64) {
	// Map each task to its next release time.
	nextRelease := make(map[int]float64)
	for _, t := range taskSet {
		nextRelease[t.ID] = 0
	}

	var runningJob *Job
	mode, dt := Normal, 0.001
	readyQueue := make([]*Job, 0)
	jobCounter, runningJobInCS := 0, false

	for currentTime := 0.0; currentTime < simulateTime; currentTime += dt {
		// Release new jobs if their release time has arrived.
		for _, t := range taskSet {
			if currentTime >= nextRelease[t.ID] {
				// In Overrun mode, only release HC tasks.
				if mode == Normal || (mode == Overrun && t.Criticality == tasks.HC) {
					jobCounter++
					newJob := &Job{
						Task:             t,
						JobID:            jobCounter,
						ReleaseTime:      nextRelease[t.ID],
						AbsoluteDeadline: nextRelease[t.ID] + t.Deadline,
						RemainingTime:    t.WCET1,
						ExecTime:         0,
					}
					// For HC tasks, choose WCET1 in Normal mode and WCET1+WCET2 in Overrun mode.
					if t.Criticality == tasks.HC && mode == Overrun {
						newJob.RemainingTime += t.WCET2
					}

					readyQueue = append(readyQueue, newJob)
					fmt.Printf("Time %.3f: Released Job %d (Task %d, Deadline=%.3f, WCET=%.3f) [Mode: %v]\n",
						currentTime, newJob.JobID, t.ID, newJob.AbsoluteDeadline, newJob.RemainingTime, mode)
				}

				// Schedule the next release for the task.
				nextRelease[t.ID] += t.Period
			}
		}

		// Scheduler decision: select a job to run.
		if len(readyQueue) > 0 {
			sort.Slice(readyQueue, func(i, j int) bool {
				return readyQueue[i].effectivePriority() < readyQueue[j].effectivePriority()
			})

			if runningJob == nil {
				// Pick the job with the smallest effective priority.
				runningJob, readyQueue = readyQueue[0], readyQueue[1:]
				fmt.Printf("Time %.3f: Starting Job %d (Task %d) with Deadline=%.3f, EffectivePriority=%.3f\n",
					currentTime, runningJob.JobID, runningJob.Task.ID, runningJob.AbsoluteDeadline, runningJob.effectivePriority())

				if cs := runningJob.getActiveCriticalSection(); cs != nil {
					fmt.Printf("Time %.3f: Job %d ENTERS critical section on Resource %d (CS: Start=%.3f, Duration=%.3f)\n",
						currentTime, runningJob.JobID, cs.ResourceID, cs.Start, cs.Duration)
				}
			} else {
				// Check if a waiting job has a lower effective priority.
				candidate := readyQueue[0]
				if candidate.effectivePriority() < runningJob.effectivePriority() {
					// If runningJob is in a critical section, allow preemption only if candidate beats its preemption level.
					if runningJob.getActiveCriticalSection() != nil {
						if candidate.effectivePriority() < float64(runningJob.Task.PreemptionLevel) {
							fmt.Printf("Time %.3f: Preempting Job %d (Task %d, EffPri=%.3f, in CS) with Job %d (Task %d, EffPri=%.3f)\n",
								currentTime, runningJob.JobID, runningJob.Task.ID, runningJob.effectivePriority(),
								candidate.JobID, candidate.Task.ID, candidate.effectivePriority())

							readyQueue = append(readyQueue, runningJob)
							runningJob, readyQueue = candidate, readyQueue[1:]
							if cs := runningJob.getActiveCriticalSection(); cs != nil {
								fmt.Printf("Time %.3f: Job %d ENTERS critical section on Resource %d (CS: Start=%.3f, Duration=%.3f)\n",
									currentTime, runningJob.JobID, cs.ResourceID, cs.Start, cs.Duration)
							}
						}
					} else {
						fmt.Printf("Time %.3f: Preempting Job %d (Task %d, EffPri=%.3f) with Job %d (Task %d, EffPri=%.3f)\n",
							currentTime, runningJob.JobID, runningJob.Task.ID, runningJob.effectivePriority(),
							candidate.JobID, candidate.Task.ID, candidate.effectivePriority())

						runningJob, readyQueue = candidate, readyQueue[1:]
						if cs := runningJob.getActiveCriticalSection(); cs != nil {
							fmt.Printf("Time %.3f: Job %d ENTERS critical section on Resource %d (CS: Start=%.3f, Duration=%.3f)\n",
								currentTime, runningJob.JobID, cs.ResourceID, cs.Start, cs.Duration)
						}
					}
				}
			}
		}

		// Check and log critical section entry/exit transitions.
		if runningJob != nil {
			cs := runningJob.getActiveCriticalSection()
			if !runningJobInCS && cs != nil {
				fmt.Printf("Time %.3f: Job %d (Task %d) ENTERS critical section on Resource %d (CS: Start=%.3f, Duration=%.3f)\n",
					currentTime, runningJob.JobID, runningJob.Task.ID, cs.ResourceID, cs.Start, cs.Duration)
			} else if runningJobInCS && cs == nil {
				fmt.Printf("Time %.3f: Job %d (Task %d) EXITS critical section\n", currentTime, runningJob.JobID, runningJob.Task.ID)
			}
			runningJobInCS = (cs != nil)
		}

		// Execute the running job for one time step.
		if runningJob != nil {
			runningJob.ExecTime += dt
			runningJob.RemainingTime -= dt

			// Check if an HC job overruns its normal (WCET1) execution in Normal mode.
			if runningJob.Task.Criticality == tasks.HC && mode == Normal && runningJob.ExecTime > runningJob.Task.WCET1 {
				mode = Overrun
				fmt.Printf("Time %.3f: Mode switch to OVERRUN triggered by Job %d (Task %d) [ExecTime=%.3f, WCET1=%.3f]\n",
					currentTime, runningJob.JobID, runningJob.Task.ID, runningJob.ExecTime, runningJob.Task.WCET1)

				// Drop pending LC jobs.
				readyQueue = dropLCJobs(readyQueue, currentTime)

				// Extend the remaining time of the HC jobs to include WCET2.
				extendRemainingTime(readyQueue)
			}

			// Job completion.
			if runningJob.RemainingTime <= 0 {
				fmt.Printf("Time %.3f: COMPLETED Job %d (Task %d) [FinishTime=%.3f, Total ExecTime=%.3f]\n",
					currentTime, runningJob.JobID, runningJob.Task.ID, currentTime, runningJob.ExecTime)
				runningJob = nil
			}
		}
	}
}
