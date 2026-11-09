package tasks

import (
	"errors"
	"fmt"
)

var ErrTaskNotFound = errors.New("Task not found")

type Priority int

// Las propiedades deben ser privadas para que no se puedan acceder desde fuera de la estructura y no se puedan modificar desde fuera de la estructura.
// las modificaciones deben ser solo mediante metodos.
type Task struct {
	id       int
	name     string
	done     bool
	priority Priority
}

func (p Priority) Validate() bool {
	return p >= 1 && p <= 3
}

func (t Task) Describe() string {
	status := "pending"

	if t.done {
		status = "done"
	}
	return fmt.Sprintf("%d [%s] %s (priority %d)", t.id, status, t.name, t.priority)

}

func (t Task) IsUrgent() bool {
	return t.priority >= 3 && !t.done
}

func AddTask(tasks map[int]Task, nextID int, name string, priority Priority) {
	tasks[nextID] = Task{id: nextID, name: name, priority: priority}

}

func MarkTaskDone(tasks map[int]Task, id int) error {
	task, exists := tasks[id]
	if !exists {
		return fmt.Errorf("task %d: %w", id, ErrTaskNotFound)
	}
	task.done = true
	tasks[id] = task
	return nil
}

func DeleteTask(tasks map[int]Task, id int) error {
	_, exist := tasks[id]
	if !exist {
		return fmt.Errorf("task %d: %w", id, ErrTaskNotFound)
	}
	delete(tasks, id)
	return nil
}
