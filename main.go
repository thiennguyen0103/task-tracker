package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type statusType string

const (
	todo       statusType = "todo"
	inProgress statusType = "in-progress"
	done       statusType = "done"
)

type commandType string

const (
	add            commandType = "add"
	update         commandType = "update"
	delete         commandType = "delete"
	markInProgress commandType = "mark-in-progress"
	markDone       commandType = "mark-done"
	list           commandType = "list"
)

const fileName = "tasks.json"

type Task struct {
	ID          int        `json:"id"`
	Description string     `json:"description"`
	Status      statusType `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: task-cli [command] [arguments]")
		return
	}

	command := commandType(args[1])

	handleTasks(command, args)
}

func loadTasks() []Task {
	var tasks []Task
	data, err := os.ReadFile(fileName)
	if err == nil {
		json.Unmarshal(data, &tasks)
	}

	return tasks
}

func handleTasks(command commandType, args []string) {
	tasks := loadTasks()

	switch command {
	case add:
		if len(args) < 3 {
			fmt.Println("Usage: task-cli add \"task name\"")
			return
		}
		description := strings.Join(args[2:], " ")
		addTask(tasks, description)
	case update:
		if len(args) < 4 {
			fmt.Println("Usage: task-cli update [id] \"description\"")
			return
		}
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid ID")
			return
		}
		description := strings.Join(args[3:], " ")
		updateTask(tasks, id, description)
	case delete:
		if len(args) < 3 {
			fmt.Println("Usage: task-cli delete [id]")
			return
		}
		id, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println("Invalid ID")
			return
		}
		deleteTask(tasks, id)
	case markInProgress:
		changeTaskStatus(tasks, args, inProgress)
	case markDone:
		changeTaskStatus(tasks, args, done)
	case list:
		if len(args) == 3 {
			listTasksByStatus(tasks, statusType(args[2]))
		} else {
			listTasks(tasks)
		}
	default:
		fmt.Println("Unknown command:", command)
	}
}

func addTask(tasks []Task, description string) {
	newID := 1
	for _, t := range tasks {
		if t.ID >= newID {
			newID += t.ID + 1
		}
	}

	now := time.Now()
	task := Task{newID, description, todo, now, now}
	tasks = append(tasks, task)
	saveTasks(tasks)
	fmt.Printf("Task added successfully (ID: %d)\n", newID)
}

func updateTask(tasks []Task, id int, description string) {
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			saveTasks(tasks)
			fmt.Printf("Task updated successfully (ID: %d)\n", id)
			return
		}
	}

	fmt.Println("Task not found")
}

func deleteTask(tasks []Task, id int) {
	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[0:i], tasks[i+1:]...)
			saveTasks(tasks)
			fmt.Printf("Task deleted successfully (ID: %d)\n", id)
			return
		}
	}

	fmt.Println("Task not found")
}

func changeTaskStatus(tasks []Task, args []string, status statusType) {
	if len(args) < 3 {
		fmt.Printf("Usage: task-cli mark-%s [id]\n", status)
		return
	}

	id, _ := strconv.Atoi(args[2])

	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			saveTasks(tasks)
			fmt.Printf("Task %d marked as %s.\n", id, status)
			return
		}
	}
	fmt.Println("Task not found.")
}

func listTasks(tasks []Task) {
	if len(tasks) == 0 {
		fmt.Println("No tasks.")
		return
	}
	for _, task := range tasks {
		printTask(task)
	}
}

func listTasksByStatus(tasks []Task, status statusType) {
	found := false
	for _, t := range tasks {
		if t.Status == status {
			printTask(t)
			found = true
		}
	}
	if !found {
		fmt.Printf("No tasks with status \"%s\".\n", status)
	}
}

func saveTasks(tasks []Task) {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	os.WriteFile(fileName, data, 0644)
}

func printTask(t Task) {
	fmt.Printf("ID: %d\nDescription: %s\nStatus: %s\nCreatedAt: %s\nUpdatedAt: %s\n\n",
		t.ID, t.Description, t.Status, t.CreatedAt.Format(time.RFC3339), t.UpdatedAt.Format(time.RFC3339))
}
