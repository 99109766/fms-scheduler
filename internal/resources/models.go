package resources

import "fmt"

type Resource struct {
	ID            int
	AssignedTasks []int
	Ceiling       int
}

func (r Resource) String() string {
	return fmt.Sprintf("Resource %d -> Tasks: %v, Ceiling: %d", r.ID, r.AssignedTasks, r.Ceiling)
}
