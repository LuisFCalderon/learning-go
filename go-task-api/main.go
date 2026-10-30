package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

type Task struct {
	ID       int
	Name     string
	Done     bool
	Priority int
}

func (t Task) Describe() string {
	status := "pending"

	if t.Done {
		status = "done"
	}
	return fmt.Sprintf("%d [%s] %s (priority %d)", t.ID, status, t.Name, t.Priority)

}

func (t Task) IsUrgent() bool {
	return t.Priority >= 3 && !t.Done
}

func findTaskIndex(tasks []Task, id int) (int, error) {
	for i, t := range tasks {
		if t.ID == id {
			return i, nil
		}
	}

	return -1, fmt.Errorf("No task was found with id: %d", id)
}

func addTask(tasks []Task, nextId int, name string, priority int) []Task {
	newTask := Task{ID: nextId, Name: name, Priority: priority}
	return append(tasks, newTask)
}

func markTaskDone(tasks []Task, id int) ([]Task, error) {
	index, error := findTaskIndex(tasks, id)
	if error != nil {
		return tasks, error
	}

	tasks[index].Done = true
	return tasks, nil

}

func deleteTask(tasks []Task, id int) ([]Task, error) {

	index, error := findTaskIndex(tasks, id)
	if error != nil {
		return tasks, error
	}

	tasks = append(tasks[:index], tasks[index+1:]...)
	return tasks, nil

}

func readLine(prompt string) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func printMenu() {
	fmt.Println("\n--- Task Manager ---")
	fmt.Println("1. Add task")
	fmt.Println("2. List tasks")
	fmt.Println("3. Mark task as done")
	fmt.Println("4. Delete task")
	fmt.Println("5. Exit")
}

func main() {
	var tasks []Task
	nextID := 1

	for {
		printMenu()
		choice := readLine("Choose an option: ")

		switch choice {
		case "1":
			name := readLine("Task name: ")
			priorityText := readLine("Priority (1-3): ")

			priority, err := strconv.Atoi(priorityText)
			if err != nil {
				fmt.Println("Priority must be a number")
				continue
			}
			tasks = addTask(tasks, nextID, name, priority)
			nextID++

			fmt.Printf("Task (%s) added!\n", name)

		case "2":
			if len(tasks) == 0 {
				fmt.Println("No task yet")
				continue
			}

			for _, t := range tasks {
				fmt.Println(t.Describe())
				if t.IsUrgent() {
					fmt.Println("this one is urgent")
				}
			}

		case "3":
			idText := readLine("Task ID to mark as done: ")
			id, err := strconv.Atoi(idText)
			if err != nil {
				fmt.Println("This is not a valid ID")
				continue
			}
			tasks, err = markTaskDone(tasks, id)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Task marked as done.")

		case "4":
			idText := readLine("Task ID to delete: ")
			id, err := strconv.Atoi(idText)
			if err != nil {
				fmt.Println("This is not a valid ID")
				continue
			}

			tasks, err = deleteTask(tasks, id)

			if err != nil {
				fmt.Println("This is not a valid ID")
				continue
			}

		case "5":
			fmt.Println("Goodbye.")
			return
		default:
			fmt.Println("Invalid option.")
		}

	}
}
