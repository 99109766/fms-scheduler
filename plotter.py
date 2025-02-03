import json
import matplotlib.pyplot as plt

def plot_periodic_schedule(schedule_file, tasks_file):
    # Load JSON data
    with open(schedule_file, 'r') as file:
        schedule = json.load(file)
    
    with open(tasks_file, 'r') as file:
        tasks = json.load(file)
    
    num_tasks = len(tasks)
    fig, axes = plt.subplots(num_tasks, 1, figsize=(12, 6), sharex=True)
    
    colors = {}
    max_t = max([task['end_time'] for task in schedule])
    for task in tasks:
        task_id = task['id']
        period = task['period']
        deadline = task['deadline']
        
        # Assign each task a unique color
        if task_id not in colors:
            colors[task_id] = plt.cm.Paired(len(colors) % 12 / 12)
        
        ax = axes[int(task_id) - 1] if num_tasks > 1 else axes

        # Draw arrows for period start and deadline
        t = 0
        while t < max_t:
            ax.annotate(r'$\uparrow$', xy=(t - 0.2, -0.1), ha='center', fontsize=12)
            if t + deadline < max_t:
                ax.annotate(r'$\downarrow$', xy=(t + deadline - 0.2, -0.1), ha='center', fontsize=12)
            t += period
        
        # Filter execution times from schedule
        task_executions = [t for t in schedule if t['task_id'] == task_id]
        
        for exec_instance in task_executions:
            start_time = exec_instance['start_time']
            end_time = exec_instance['end_time']
            ax.barh(0, end_time - start_time, left=start_time, height=0.5, color=colors[task_id], edgecolor='black')
        
        ax.set_yticks([])
        ax.set_ylabel(f'Task {task_id}', rotation=0, labelpad=20, va='center')
    
    # Formatting
    plt.xlabel("Time")
    fig.suptitle("Periodic Task Scheduling Charts")
    plt.show()

def plot_qoc_schedule(schedule_file):
    # Load schedule data
    with open(schedule_file, "r") as file:
        schedule = json.load(file)

    # Compute response times
    response_times = []
    for task in schedule:
        response_time = task["end_time"] - task["start_time"]
        response_times.append(response_time)

    # Plot response time distribution
    plt.figure(figsize=(10, 5))
    plt.hist(response_times, bins=20, edgecolor="black")
    plt.xlabel("Response Time")
    plt.ylabel("Frequency")
    plt.title("Response Time Distribution")
    plt.grid(axis="y", linestyle="--", alpha=0.7)
    plt.show()

plot_qoc_schedule("schedule.json")
plot_periodic_schedule("schedule.json", "tasks.json")