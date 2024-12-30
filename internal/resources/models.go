package resources

import "fmt"

type Resource struct {
	ID            int
	AssignedTasks []int
}

func (r Resource) String() string {
	return fmt.Sprintf("Resource %d -> %v", r.ID, r.AssignedTasks)
}
