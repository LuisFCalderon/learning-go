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

func main() {

	var task Task

	fmt.Printf("ID: %d, Name: %s, Done: %t, Priority: %d\n", task.ID, task.Name, task.Done, task.Priority)
	fmt.Println("ID: ", task.ID)             // (no aceptable) No lo considero aceptable debido a que los ID incrementales rara vez inicializa en 1 y no en 0.
	fmt.Println("Name: ", task.Name)         //(no aceptable) Si es para una tarea no tiene mucho sentido porque no existe una tarea sin nombre.
	fmt.Println("Done: ", task.Done)         //(aceptable) Considero que este es correcto porque una tarea al iniciar no esta terminada.
	fmt.Println("Priority: ", task.Priority) //(aceptable?) Puede o no ser correcto aunnque cuando se crea una tarea le asignas un prioridad siempre, por lo que que inicie en 0 es incosiste por lo menos en mi opinion.

	fmt.Print("Name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Print("Priority: ")
	priorityText, _ := reader.ReadString('\n')
	priorityText = strings.TrimSpace(priorityText)
	priority, err := strconv.Atoi(priorityText)

	if err != nil || (priority < 1 || priority > 3) {
		fmt.Println("Priority must be a number between 1 and 3")
		return
	}

	task = Task{
		ID:       1,
		Name:     name,
		Priority: priority,
	}

	fmt.Printf("ID: %d, Name: %s, Done: %t, Priority: %d\n", task.ID, task.Name, task.Done, task.Priority)
	//por qué validar el int y no el texto crudo? Bajo mi concepto porque es mas simple evaluar un entero entre rangos que la forma del texto es decir, yo doy un rango entre 1 y 3 pero para el texto presiento que un usuario si le doy campo abierto escribiria alto,medio,bajo por ejemplo y no un rango de de numeros y fallaria todo el flujo porque esta buscando un posible caracter que no se encontraía.

}
