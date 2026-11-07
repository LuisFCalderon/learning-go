package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

type Priority int

type Task struct {
	ID       int
	Name     string
	Done     bool
	Priority Priority
}

func (p Priority) Validate() bool {
	return p >= 1 && p <= 3
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

func addTask(tasks []Task, nextId int, name string, priority Priority) []Task {
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

//¿Un tipo con nombre, type Priority int, cierra la grieta de que priority se puede pasar como int?
//Si, porque podemos usar metodos para validar el rango que buscanmos y evitar que se pase un valor que no sea entre 1 y 3.

//¿para qué sirve darle nombre a un int?
//Para poder usar extender el tipo int con metodos y funciones, considero que se puede hacer mas robusto darle metodos y funciones a un tipo que ya existe, lo que lo hace mas testiable tambien.

func main() {
	var tasks []Task
	nextID := 1

	for {
		printMenu()
		choice := readLine("Choose an option: ")

		switch choice {
		case "1":
			name := readLine("Task name: ")
			priorityInput := readLine("Priority (1-3): ")

			priorityInt, err := strconv.Atoi(priorityInput)
			if err != nil {
				fmt.Println("This is not a valid priority")
				continue
			}
			if !Priority(priorityInt).Validate() {
				fmt.Println("Priority must be a number between 1 and 3")
				continue
			}
			// Convertimos el int a Priority porque el 4to parametro de addTask es de tipo Priority ademas de que se puede validar el rango de prioridad.
			//Seria en este punto hacer la conversion a Priority para asi evitamos la falla en la funcion addTask
			priority := Priority(priorityInt)

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
