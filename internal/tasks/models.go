package tasks

import (
	"fmt"
)

type CriticalityLevel int

const (
	LC CriticalityLevel = iota
	HC
)

type CriticalSection struct {
	ResourceID int     `json:"resource_id"`
	Start      float64 `json:"start"`
	Duration   float64 `json:"duration"`
}

func (cs CriticalSection) End() float64 {
	return cs.Start + cs.Duration
}

type Task struct {
	ID               int                `json:"id"`
	Criticality      CriticalityLevel   `json:"criticality"`
	Period           float64            `json:"period"`
	Deadline         float64            `json:"deadline"`
	WCET1            float64            `json:"wcet1"`
	WCET2            float64            `json:"wcet2"`
	AssignedResIDs   []int              `json:"assigned_res_ids"`
	CriticalSections []*CriticalSection `json:"critical_sections"`
	Priority         int                `json:"-"`
	PreemptionLevel  int                `json:"-"`
}

func (t *Task) Utilization() float64 {
	return t.WCET1 / t.Period
}

func (t *Task) MaxUtilization() float64 {
	if t.Criticality == LC {
		return t.WCET1 / t.Period
	}
	return (t.WCET1 + t.WCET2) / t.Period
}

func (t *Task) String() string {
	if t.Criticality == LC {
		return fmt.Sprintf(
			"[Task %d | LC | Period=%.2f | Deadline=%.2f | WCET=%.2f | Util=%.2f | Priority=%d | PreemptionLevel=%d | Res=%v]",
			t.ID, t.Period, t.Deadline, t.WCET1, t.Utilization(), t.Priority, t.PreemptionLevel, t.AssignedResIDs)
	}
	return fmt.Sprintf(
		"[Task %d | HC | Period=%.2f | Deadline=%.2f | WCET1=%.2f | WCET2=%.2f | Util=%.2f | MaxUtil=%.2f | Priority=%d | PreemptionLevel=%d | Res=%v]",
		t.ID, t.Period, t.Deadline, t.WCET1, t.WCET2, t.Utilization(), t.MaxUtilization(), t.Priority, t.PreemptionLevel, t.AssignedResIDs)
}
