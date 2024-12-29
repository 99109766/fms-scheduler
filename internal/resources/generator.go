package resources

import "fmt"

// GenerateResources creates a list of resource IDs [0..numResources-1].
func GenerateResources(numResources int) []Resource {
	resources := make([]Resource, numResources)
	for i := 0; i < numResources; i++ {
		resources[i] = Resource{ID: i, AssignedTasks: make([]int, 0)}
	}
	fmt.Printf("Generated %d Resources: %v\n", numResources, resources)
	return resources
}
