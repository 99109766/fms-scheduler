package resources

import "fmt"

type Resource struct {
	ID            int   `json:"id"`
	AssignedTasks []int `json:"assigned_tasks"`
	Ceiling       int   `json:"ceiling"`
}

func (r Resource) String() string {
	return fmt.Sprintf("Resource %d -> Tasks: %v, Ceiling: %d", r.ID, r.AssignedTasks, r.Ceiling)
}
