package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/99109766/fms-scheduler/config"
	"github.com/99109766/fms-scheduler/internal/resources"
	"github.com/99109766/fms-scheduler/internal/scheduler"
	"github.com/99109766/fms-scheduler/internal/tasks"
)

func main() {
	// Parse flags
	configPathPtr := flag.String("config", "", "Path to the configuration file (YAML format)")
	flag.Parse()

	if configPathPtr == nil {
		log.Fatal("Config file path must be provided using --config or -c flag")
	}
	configPath := *configPathPtr

	// Load Configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Generate tasks using UUnifast without any resource assignments
	taskSet := tasks.GenerateTasksUUnifast(cfg)

	fmt.Println("=== Generated Tasks ===")
	for _, t := range taskSet {
		fmt.Println(t)
	}

	// Generate resources without any assignments
	resourceList := resources.GenerateResources(cfg.NumResources)

	// Assign resources to tasks (e.g., nested resource usage)
	tasks.AssignResourcesToTasks(cfg, taskSet, resourceList)

	fmt.Println("\n=== Resource Assignments ===")
	for _, r := range resourceList {
		fmt.Printf("Resource %d assigned to tasks: %v\n", r.ID, r.AssignedTasks)
	}

	// Assign critical sections to tasks
	tasks.AssignCriticalSections(cfg, taskSet, resourceList)

	fmt.Println("\n=== Tasks and Assigned Critical Sections ===")
	for _, t := range taskSet {
		fmt.Printf("Task %d (Criticality: %v) Critical Sections:\n", t.ID, t.Criticality)
		for _, cs := range t.CriticalSections {
			fmt.Printf("  - Resource %d: %.2f\n", cs.ResourceID, cs.Duration)
		}
	}

	// Determine priority levels or preemption levels
	scheduler.DeterminePriorityLevels(taskSet)

	fmt.Println("\n=== Final Task Details (with assigned priorities) ===")
	for _, t := range taskSet {
		fmt.Println(t)
	}
}
