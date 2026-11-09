# Curso de Go: cuándo elegir un tipo

Curso personal guiado por Claude. El tema único es **cuándo elegir un tipo**, usando el
gestor de tareas de este repo (`go-task-api/main.go`) como ejemplo. Una lección por
sesión; no se avanza hasta que yo escribo el ejercicio y se revisa. Sin Gin, sin base de
datos, sin Docker, sin concurrencia.

**Reglas de revisión:** primero si elegí el tipo correcto, después si compila, después si
hace lo pedido. Si el tipo está mal, se me dice qué caso real se rompe, sin reescribirme
la solución.

## Progreso

| # | Lección | Estado |
|---|---------|--------|
| 1 | El valor cero | ✅ Aprobada (2025-11-03) |
| 2 | `int`, `string` y `bool`; la frontera de `strconv.Atoi` | ✅ Aprobada (2026-11-04) |
| 3 | Tipo con nombre: `type Priority int` | ✅ Aprobada (2026-11-07) |
| 4 | Slice o map | ✅ Aprobada (2026-11-07) |
| 5 | Devolver la colección o no (semántica de referencia) | ✅ Aprobada (2026-11-07) |
| 6 | Puntero o valor | ✅ Aprobada (2026-11-07) |
| 7 | El tipo `error` | ✅ Aprobada (2026-11-08) |
| 8 | Paquetes | ✅ Aprobada (2026-11-08) |
| 9 | Tests de tabla | Pendiente |
| 10 | La base del proyecto (`internal/config`, `internal/domain`, `cmd/`) | Pendiente |

## ¿Y después de la lección 10?

Este curso termina a propósito **antes** de lo avanzado: su única pregunta es *cuándo
elegir un tipo*, y concurrencia, HTTP, Gin, PostgreSQL y Docker quedan excluidos por
diseño. Al aprobar la lección 10 (la base del proyecto), se puede armar un segundo curso
sobre esa misma base, en este orden natural:

1. **Concurrencia**: goroutines, `channels`, `sync.WaitGroup`, `context.Context` — y la
   misma pregunta de siempre aplicada ahí (¿channel o mutex?, ¿`context` o parámetro?).
2. **HTTP**: `net/http` de la librería estándar antes que Gin, sirviendo el paquete
   `tasks` ya escrito.
3. **Persistencia**: PostgreSQL detrás de una interfaz definida en `internal/domain`.

Nada de esto se adelanta: primero las 10 lecciones.

## Lo aprendido

### Lección 1 — El valor cero

- Todo tipo tiene valor cero garantizado: `int` → `0`, `string` → `""`, `bool` → `false`,
  slice y puntero → `nil`. Un struct nace con cada campo en su valor cero.
- La pregunta de diseño: ¿el valor cero es un estado válido del dominio o un estado roto?
  En `Task`: `Done == false` es válido (recién creada = no hecha); `ID == 0` y
  `Name == ""` no son estados de ninguna tarea real.
- Que `Priority == 0` sea inválido es *útil*: permite detectar "nadie asignó prioridad".
- Para inspeccionar strings vacíos, `%q` (imprime `""` con comillas).

### Lección 2 — `int`, `string`, `bool` y la frontera

- `string` solo mientras el dato es texto crudo sin interpretar. `int` cuando se calcula
  o compara como número. `bool` cuando solo hay dos estados.
- El teclado siempre entrega `string`. `strconv.Atoi` es **la frontera**: se cruza una
  vez, lo antes posible, y del otro lado nunca se vuelve a tocar el texto.
- `Atoi` solo garantiza "es un número"; el rango (1–3) se valida sobre el `int`, no sobre
  el texto: `"02"` significa 2 y `Atoi` lo entiende; comparar strings compara dibujos,
  no significados.
- Bug propio del que aprendí: `priority < 1 && priority > 3` es el conjunto vacío (ningún
  número cumple ambas); "fuera de rango" es una unión, se escribe con `||`.
- Frase clave del curso: **el tipo correcto no es el que representa los estados válidos,
  es el que hace irrepresentables los inválidos** (`t.Done = 7` no compila con `bool`).
- Límite de la lección: la validación en `main` es una costumbre, no una garantía;
  `addTask(..., 99)` sigue compilando. Eso lo resuelve la lección 3.

### Lección 3 — Tipo con nombre: `type Priority int`

- `type Priority int` crea un tipo *distinto* de `int`, no un alias. Sirve para dos cosas:
  dar significado y permitir colgarle métodos (la regla de validación vive pegada al tipo,
  en un solo sitio testeable, no repartida por `main`).
- **Lo que el tipo SÍ cierra:** pasar un `int` *variable* donde se espera `Priority` no
  compila (`cannot use n (variable of type int) as Priority value`). Eso elimina los bugs
  de mezcla: meter `nextID` en el hueco de la prioridad ya no pasa el compilador.
- **Lo que el tipo NO cierra:** `addTask(..., 99)` con un *literal* sigue compilando. Un
  literal es una constante sin tipo y Go la convierte sola a `Priority`. Go **no tiene
  tipos con rango**: el tipo con nombre no hace irrepresentable el 99.
- Quien cierra la grieta del 99 no es el tipo, es la *llamada* a `Valid()` — y eso depende
  de que el que llama se acuerde, igual que en la lección 2. La ganancia real: la regla
  está centralizada y es testeable, no que el inválido sea imposible.
- Convención idiomática: un predicado se llama `Valid()` / `IsValid()` (exportado), no
  `validate()` en minúscula; importará cuando esto sea el paquete `tasks` (lección 10).
- La conversión en la frontera (`Priority(Atoi(...))`) va *después* de confirmar que el
  texto era número; convertir basura antes no rompe aquí, pero es al revés.

### Lección 4 — Slice o map

- **Slice** = secuencia ordenada; bueno para recorrer en orden y acceder por posición.
  Buscar por una clave (ej. ID) obliga a un barrido lineal O(n).
- **Map** (`map[int]Task`) = asociación clave→valor; búsqueda, escritura y borrado por
  clave en O(1). `delete(m, k)` sobre una clave ausente no da error (no-op).
- **Gotcha de Go:** el orden de recorrido de un map es *aleatorio a propósito* —
  el runtime lo mezcla en cada ejecución para que nadie dependa del orden por casualidad.
  Precio del map aquí: se pierde el orden de inserción al listar.
- Los valores de un map **no son direccionables**: `m[k].Campo = x` no compila. Patrón
  idiomático: sacar la copia, modificarla y reescribir (`t := m[k]; t.Done = true; m[k] = t`).
- Detectar "no existe" con el segundo valor de retorno: `v, ok := m[k]`. Base del
  not-found idiomático, en vez de un barrido que devuelve `-1`.
- Decisión para esta app: map, por búsqueda/borrado por ID y escalabilidad; el orden no
  es parte del dominio (se podría ordenar por prioridad al listar si hiciera falta).

### Lección 5 — Devolver la colección o no (semántica de referencia)

- Al pasar un argumento, Go **siempre copia**. Lo que cambia es *qué* se copia.
- Un **slice** es una cabecera (puntero al arreglo, longitud, capacidad). `append` puede
  realocar (nueva dirección → cabecera nueva) o cambiar la longitud; en ambos casos la
  cabecera de `main` queda desfasada, así que **hay que devolver el slice**.
- Matiz importante: **no todo cambio a un slice obliga a devolver**. Editar un elemento en
  su sitio (`s[i].Campo = x`) escribe en el arreglo compartido y el llamador lo ve sin
  retorno. Solo la *cabecera* (realocación o longitud) obliga a devolver.
- Un **map** también copia su handle al pasarlo, pero el handle apunta a la **misma** tabla
  hash; no hay longitud/capacidad que recolocar, así que `append`-equivalente no existe y
  mutar (set/delete/reescribir clave) se ve en el llamador **sin devolver nada**.
- Por eso las funciones del map no necesitan devolver la colección: `addTask` sin retorno,
  `markTaskDone`/`deleteTask` devuelven solo `error`.
- Frase precisa: *se copia la cabecera/handle, nunca los datos a los que apunta; con el
  slice hay que devolver cuando cambia la cabecera, con el map nunca.*

### Lección 6 — Puntero o valor

- Al pasar un `struct` (o llamar un método con receptor por valor), Go **copia el struct
  entero**. Receptor por valor = copia; receptor por puntero (`*Task`) = el original.
- **Valor** cuando el método solo *lee* y el struct es chico (`Describe`, `IsUrgent`).
  **Puntero** cuando *muta* el original, o el struct es grande y copiar pesa.
- El receptor es documentación: `func (t *Task) Describe()` le miente al lector — anuncia
  mutación que no ocurre.
- Muro clave (ata con lección 4): un método con receptor puntero **no se puede llamar
  sobre un valor no direccionable** — `tasks[id].MarkDone()` no compila
  (`cannot call pointer method ... on Task`). El map puede mover/rehashear sus entradas,
  así que Go prohíbe `&tasks[id]`; sin dirección estable, no hay receptor puntero.
- Por eso en un map el patrón para mutar es el copy-back (`t := m[k]; t.Done = true;
  m[k] = t`), no un método puntero.
- Un receptor por valor se puede llamar sobre *cualquier cosa* (map, literal, copia);
  ponerle puntero a un lector lo vuelve frágil sin ganar nada.
- Estilo Go: ser consistente — no mezclar receptores valor y puntero en un mismo tipo sin
  razón.

### Lección 7 — El tipo `error`

- `error` es una **interfaz**: `type error interface { Error() string }`. Cualquier tipo con
  `Error() string` es un `error`; por eso el conjunto de errores es abierto y cada paquete
  define los suyos (con campos/métodos propios) y el llamador los maneja a todos igual.
- El fallo es un **valor que se devuelve**, no una excepción que se lanza. Go no usa
  try/catch para esto.
- **`error`** → fallo esperado, contemplado, recuperable (no existe, no parsea, no está).
  **`panic`** → algo que nunca debería pasar: bug del programador, invariante roto, estado
  imposible. El ID que el usuario tecleó mal es `error`, no `panic`.
- Comparar el *texto* de un error es frágil (mismo error que la lección 2: comparar dibujos,
  no significados). Se pregunta por **identidad**, no por texto.
- Patrón: **error centinela** declarado una vez (`var ErrTaskNotFound = errors.New(...)`),
  devuelto envuelto con `%w` (`fmt.Errorf("id %d: %w", id, ErrTaskNotFound)`), y reconocido
  con `errors.Is(err, ErrTaskNotFound)`. `%w` conserva el mensaje legible *y* deja que
  `errors.Is` lo detecte a través de la envoltura.
- Cada error se chequea donde nace y según su naturaleza: el de `strconv.Atoi` (parseo) con
  `err != nil`; el del dominio (`markTaskDone`/`deleteTask`) con `errors.Is`. No cruzarlos.
- Convención: el string de un error va en minúscula y sin punto final (se concatena al
  envolverlo).

### Lección 8 — Paquetes

- Un paquete es una **cápsula**: la frontera se decide con una sola cosa, la **mayúscula
  inicial**. Mayúscula = exportado (parte del contrato público); minúscula = interno (nadie
  de afuera puede ni escribir su nombre; el compilador lo hace irrepresentable).
- Dominio (`Task`, `Priority`, `Validate`, `Describe`, `IsUrgent`, `AddTask`,
  `MarkTaskDone`, `DeleteTask`, `ErrTaskNotFound`) → paquete `tasks/`. Consola (`reader`,
  `readLine`, `printMenu`, el loop/`switch`) → `main`. El `import` se arma con el nombre del
  módulo en `go.mod` + la carpeta (`taskapi/tasks`), y se usa como `tasks.Task`, etc.
- Grieta corregida: **público no es "porque yo lo uso desde afuera"; es "porque el de afuera
  DEBE poder tocarlo"**. El contraejemplo son los campos de `Task`: están en mayúscula pero
  `main` nunca toca un campo directamente (solo métodos y funciones), así que podrían ser
  privados y el programa seguiría compilando.
- Qué se rompe con campos públicos: cualquiera puede construir `tasks.Task{Priority: 99}` o
  hacer `t.Done = true` **saltándose** `Validate()` y las puertas del paquete. Campos
  privados obligan a entrar por `AddTask`/métodos, que es donde vive la validación de las
  lecciones 2/3/7. Eso es encapsulación de verdad, no solo separar archivos.
- **Corregido en la misma sesión:** los cuatro campos de `Task` pasaron a privados
  (`id, name, done, priority`); ahora el único camino para crear/mutar una tarea es
  `AddTask`/`MarkTaskDone`. También se limpió `ErrTaskNotFound` a minúscula sin punto.
- Asimetría a recordar: campos privados tampoco se leen desde afuera. Si `main` necesitara un
  dato crudo, se expone un **método lector** (`func (t Task) Name() string`), nunca se vuelve
  el campo público.
