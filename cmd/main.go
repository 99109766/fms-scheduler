package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/99109766/fms-scheduler/config"
	"github.com/99109766/fms-scheduler/internal/resources"
	"github.com/99109766/fms-scheduler/internal/scheduler"
	"github.com/99109766/fms-scheduler/internal/tasks"
)

func main() {
	// Parse flags
	configPathPtr := flag.String("config", "", "Path to the configuration file (YAML format)")
	flag.Parse()

	if configPathPtr == nil || *configPathPtr == "" {
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

	resourceList := resources.GenerateResources(cfg.NumResources)
	tasks.AssignResourcesToTasks(cfg, taskSet, resourceList)

	fmt.Println("\n=== Resource Assignments ===")
	for _, r := range resourceList {
		fmt.Printf("Resource %d assigned to tasks: %v\n", r.ID, r.AssignedTasks)
	}

	tasks.AssignCriticalSections(cfg, taskSet, resourceList)

	fmt.Println("\n=== Tasks and Assigned Critical Sections ===")
	for _, t := range taskSet {
		fmt.Printf("Task %d (Criticality: %v) Critical Sections:\n", t.ID, t.Criticality)
		for _, cs := range t.CriticalSections {
			fmt.Printf("  - Resource %d: Start=%.2f, Duration=%.2f, End=%.2f\n",
				cs.ResourceID, cs.Start, cs.Duration, cs.Start+cs.Duration)
		}
	}

	tasks.DeterminePriorityLevels(taskSet)
	tasks.ComputeResourceCeilings(taskSet, resourceList)
	tasks.AssignPreemptionLevels(taskSet, resourceList)

	fmt.Println("\n=== Resources with Ceilings ===")
	for _, r := range resourceList {
		fmt.Printf("Resource %d: Ceiling = %d, Assigned Tasks = %v\n", r.ID, r.Ceiling, r.AssignedTasks)
	}

	fmt.Println("\n=== Tasks with Preemption Levels ===")
	for _, t := range taskSet {
		fmt.Printf("Task %d: Base Priority = %d, Preemption Level = %d\n", t.ID, t.Priority, t.PreemptionLevel)
	}

	fmt.Println("\n=== Running Scheduler Simulation ===")
	schedule, err := scheduler.RunScheduler(taskSet, cfg.SimulateTime)
	if err != nil {
		log.Fatalf("Error running scheduler: %v", err)
	}

	fmt.Println("\n=== Final Scheduler ===")
	encoded, err := json.MarshalIndent(schedule, "", "  ")
	if err != nil {
		log.Fatalf("Error encoding schedule: %v", err)
	}
	fmt.Println(string(encoded))

	file, err := os.Create("schedule.json")
	if err != nil {
		log.Fatalf("Error creating schedule file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(encoded)
	if err != nil {
		log.Fatalf("Error writing schedule file: %v", err)
	}

	fmt.Println("\n=== Schedule written to schedule.json ===")

	fmt.Println("\n=== Done ===")
}
