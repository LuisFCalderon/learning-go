package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"taskapi/tasks"
)

var reader = bufio.NewReader(os.Stdin)

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


**************************Excercise 6*****************************
Error es una interfaz porque podemos extenderla tanto como queramos, porque al final como programadores definimos los errores que tenemos de acuerdo a las funciones y metodos que manejamos.
Asumo que por eso fue pensada asi porque podemos tener tantos errores como funciones o metodos.
 ¿cuándo devolverías un error y cuándo harías panic?
 Panic cuando algo que nunca debe pasar.
 error para un fallo esperado,contemplado y recuperable

*/

func main() {
	/* var tasks []Task */
	taskStore := map[int]tasks.Task{}

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
			if !tasks.Priority(priorityInt).Validate() {
				fmt.Println("Priority must be a number between 1 and 3")
				continue
			}
			// Convertimos el int a Priority porque el 4to parametro de addTask es de tipo Priority ademas de que se puede validar el rango de prioridad.
			//Seria en este punto hacer la conversion a Priority para asi evitamos la falla en la funcion addTask
			priority := tasks.Priority(priorityInt)

			tasks.AddTask(taskStore, nextID, name, priority)
			nextID++
			fmt.Printf("Task (%s) added!\n", name)
		case "2":
			if len(taskStore) == 0 {
				fmt.Println("No task yet")
				continue
			}

			for _, t := range taskStore {
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
			err = tasks.MarkTaskDone(taskStore, id)
			if errors.Is(err, tasks.ErrTaskNotFound) {
				fmt.Println("Task not found")
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

			err = tasks.DeleteTask(taskStore, id)

			if errors.Is(err, tasks.ErrTaskNotFound) {
				fmt.Println("Task not found")
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
