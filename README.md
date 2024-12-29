# fms-scheduler

#### Project No. 9: Resource Management in Mixed-Criticality Systems - Flight Management

## Overview
The flight management system is a mixed-criticality system that includes tasks with high criticality, such as collision avoidance, and tasks with low criticality, such as cabin temperature control. Failure to correctly execute low-criticality tasks does not pose serious risks and only reduces the quality of service, but failure to correctly execute high-criticality tasks can lead to catastrophic outcomes.

As a result, a two-level mixed-criticality system, which is a type of mixed-criticality system, includes periodic LC and HC tasks. LC tasks have a single worst-case execution time, while HC tasks have two worst-case execution times. The system starts in a normal state, and if one of the HC tasks exceeds its smaller execution time, the system enters an overrun state, where high-criticality tasks are executed with their larger execution times.

As a real-time systems engineer, you must schedule the tasks under the **Policy Resource Stack** protocol using the **EDF-ER** scheduling algorithm, considering the multi-unit nature of shared resources. In this implementation, resources should be allocated to tasks such that both nested and serial access to resources are supported.

In addition to using the predefined set of tasks labeled as **FMS**, you should use the **Uunifast** algorithm to generate a synthetic set of tasks and report the correctness of your implementation.

## Expected Outputs:  

**Schedulability charts based on utilization rates of 0.5, 0.3, and 0.75:**
- Generate 50 tasks with an HC-to-LC task ratio of 1, each resource having a random unit count between 1 and 5, total resources being 10, and the number of critical sections in each task being a random value between 0 and 8.
- Generate 50 tasks with an HC-to-LC task ratio of 1, each resource having a random unit count between 1 and 5, total resources being 15, and the number of critical sections in each task being a random value between 6 and 10.

**Quality of Service charts based on utilization rates of 0.5, 0.7, and 0.9:**
- Generate 50 tasks with an HC-to-LC task ratio of 1, each resource having a random unit count between 1 and 5, total resources being 10, and the number of critical sections in each task being a random value between 0 and 8.
- Generate 50 tasks with an HC-to-LC task ratio of 1, each resource having a random unit count between 1 and 5, total resources being 15, and the number of critical sections in each task being a random value between 6 and 10.

---

### Phase 1:
In the first phase, for the generated tasks, implement the proposed algorithm for generating and mapping resources, allocating critical sections to tasks, and determining priority levels. The results should be reported as output.

### Phase 2:
In the final phase, all requested charts must be reported.
