package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"

	"github.com/99109766/fms-scheduler/config"
	"github.com/99109766/fms-scheduler/internal/resources"
	"github.com/99109766/fms-scheduler/internal/scheduler"
	"github.com/99109766/fms-scheduler/internal/tasks"
)

func main() {
	// --- 0) Parse flags ---
	var configPath string
	pflag.StringVar(&configPath, "config", "", "Path to the configuration file (YAML format)")
	pflag.Parse()

	if configPath == "" {
		log.Fatal("Config file path must be provided using --config flag")
	}

	// --- 1) Load Configuration ---
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// --- 2) Generate tasks using UUnifast ---
	taskSet, err := tasks.GenerateTasksUUnifast(cfg.NumTasks, cfg.TotalUtil)
	if err != nil {
		log.Fatalf("Error generating tasks: %v", err)
	}

	// Assign random criticalities (some LC, some HC)
	// and for HC tasks, assign two different WCET values:
	tasks.AssignMixedCriticality(taskSet)

	// Print out the tasks info
	fmt.Println("=== Generated Tasks ===")
	for _, t := range taskSet {
		fmt.Println(t)
	}

	// --- 3) Generate Resources & Map them to tasks ---
	resourceList := resources.GenerateResources(cfg.NumResources)

	// Assign resources to tasks (e.g., nested resource usage)
	tasks.AssignResourcesToTasks(taskSet, resourceList)

	fmt.Println("\n=== Resource Assignments ===")
	for _, r := range resourceList {
		fmt.Printf("Resource %d assigned to tasks: %v\n", r.ID, r.AssignedTasks)
	}

	// --- 4) Assign critical sections to tasks ---
	tasks.AssignCriticalSections(taskSet, resourceList)

	// --- 5) Determine priority levels or preemption levels ---
	scheduler.DeterminePriorityLevels(taskSet)

	// Print results
	fmt.Println("\n=== Final Task Details (with assigned priorities) ===")
	for _, t := range taskSet {
		fmt.Println(t)
	}
}
