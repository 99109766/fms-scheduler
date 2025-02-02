package tasks

import (
	"math"
	"math/rand"

	"github.com/99109766/fms-scheduler/config"
)

// GenerateTasksUUnifast generates a set of tasks whose sum of utilization = totalUtil.
func GenerateTasksUUnifast(cfg *config.Config) []*Task {
	numTasks, totalUtil := cfg.NumTasks, cfg.TotalUtility

	// Apply UUnifast algorithm
	utilizations := uUniFast(numTasks, totalUtil)

	// Create Task structures
	tasks := make([]*Task, numTasks)
	for i := 0; i < numTasks; i++ {
		period := cfg.PeriodRange[0] + rand.Float64()*(cfg.PeriodRange[1]-cfg.PeriodRange[0])
		wcet := utilizations[i] * period

		tasks[i] = &Task{
			ID:       i + 1,
			WCET1:    wcet,
			Period:   period,
			Deadline: period, // By default, deadline = period
		}
	}

	assignRandomCriticality(cfg, tasks)

	return tasks
}

// assignRandomCriticality assigns random criticality to tasks.
func assignRandomCriticality(cfg *config.Config, tasks []*Task) {
	for _, t := range tasks {
		if rand.Float64() < cfg.HighRatio {
			t.Criticality = HC
			t.WCET2 = (cfg.WCETRatio[0] + rand.Float64()*(cfg.WCETRatio[1]-cfg.WCETRatio[0])) * t.WCET1
		} else {
			t.Criticality = LC
			t.WCET2 = 0
		}
	}
}

// uUniFast is the internal function implementing the UUniFast algorithm.
func uUniFast(n int, U float64) []float64 {
	sumU := U
	utils := make([]float64, n)
	for i := 1; i < n; i++ {
		next := sumU * (math.Pow(rand.Float64(), 1.0/float64(n-i)))
		utils[i-1] = sumU - next
		sumU = next
	}
	utils[n-1] = sumU
	return utils
}
