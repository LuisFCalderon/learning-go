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

/* Nueva apreciacion.
Entiendo que cuando quiero modificar el task necesito la direccion, por lo que para las metodos Describe e IsUrgent no es necesario el puntero.
Dado esto seria ideal para estos metodos solo es suficiente con ser receptores por valor

/* func (t *Task, id Int) MarkDone() bool {
	t.Done = true
	return t
}
 C:\Users\calde\OneDrive\Documentos\Learning\learning-go\go-task-api>
                                                                        go run .
# taskapi
.\main.go:91:12: cannot call pointer method MarkDone on Task

Este error pasa porque los maps pueden cambiar de posicion porlo que no tenemos la referencia de la posible ubicacion en memoria del Task, debido a esto diriamos que no sabemos donde esta el original(direccion del original+)
*/

/* func findTaskIndex(tasks []Task, id int) (int, error) {
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

	append(tasks[:index], tasks[index+1:]...)
	return tasks, nil

} */

func addTask(tasks map[int]Task, nextID int, name string, priority Priority) {
	tasks[nextID] = Task{ID: nextID, Name: name, Priority: priority}

}

func markTaskDone(tasks map[int]Task, id int) error {
	task, exists := tasks[id]
	if !exists {
		return fmt.Errorf("Task with ID %d not found", id)
	}
	task.Done = true
	tasks[id] = task
	return nil
}

func deleteTask(tasks map[int]Task, id int) error {
	_, exist := tasks[id]
	if !exist {
		return fmt.Errorf("Task with ID: %d not found", id)
	}
	delete(tasks, id)
	return nil
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

//****************************************************************
/*  Exercise 4

¿El dominio necesita orden, necesita búsqueda por clave, o las dos?
Bajo mi criterio puede ser que necesite un orden, pero no es estrictamente necesario porque conceptualmente el ser humano siempre ordena las tareas por prioridad y no por como van llegando o estan dispuestas.
Creo que puden trabajar de ambas maneras con pero por escalabilidad eligiria usar un slice.
¿para esta app te quedas con el map o vuelves al slice?
Yo me quedaria con el map, porque me simplifica o me ahorra el proceso de hacer busqueda por indeces por un ID, que en una lista pequeña no se siente pero para una lista de cientos de miles o millones ya el rendimiento al recorrer todo tiende a reducirse
Corre la opción "listar" varias veces con 3 o 4 tareas. ¿Qué le pasó al orden? ¿Por qué?
El orden cambio en todas las ejecuciones, por lo que investigue un map es una tablas hash que se guardan segun este hash, Esto se hizo para que un programa que depnede de un orden fijo falle y que no funcione por casualidad.
¿para esta app te quedas con el map o vuelves al slice?

Me quedaria con el map porque me simplifica muchas cosas y escala mejor

**************** EXercise 5******************************
 ¿qué pintan los return tasks y los tasks = ...? ¿Son necesarios, o son herencia del slice que ya no hace nada?
 Efectivamente era algo hecho por los slices, debibo a lo que se habia mencionado de que los slices siempre se genera una copia del mismo y las modificaciones al final es una variacion del slice, pero es un slice nuevo lo que hace que los returns sean algo heredado
 En pocas palabras no son necesarios o no pintan nada cuando pasamos a los maps
 ¿por qué esto funciona con el map pero habría roto con el slice?
Mi interpretacion es que cada cambio sobre el slice, de cierta manera genera uno nuevo y el previo a la modificacion siempre es el que esta en flujo que se esta ejecutando por lo que al devolverlo estariamos actualizando la version que esta en el flujo previo. Por eso hay que devolverlo y si no se romperia el flujo que estas corriendo.
Ademas, entiendo que como funciona go es que con las maps, no genera compia si no un handle que permite alterar el valor del map en el flujo.

*/

func main() {
	/* var tasks []Task */
	tasks := map[int]Task{}

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

			addTask(tasks, nextID, name, priority)
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
			err = markTaskDone(tasks, id)
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

			err = deleteTask(tasks, id)

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
