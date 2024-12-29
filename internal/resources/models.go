package resources

import "fmt"

// Resource model. For demonstration, we store only an ID and which tasks are using it.
type Resource struct {
	ID            int
	AssignedTasks []int
}

func (r Resource) String() string {
	return fmt.Sprintf("Resource %d -> %v", r.ID, r.AssignedTasks)
}
