package resources

// GenerateResources creates a list of resource IDs.
func GenerateResources(numResources int) []*Resource {
	resources := make([]*Resource, numResources)
	for i := 0; i < numResources; i++ {
		resources[i] = &Resource{
			ID:            i,
			AssignedTasks: make([]int, 0),
		}
	}
	return resources
}
