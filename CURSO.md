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
| 4 | Slice o map | 🔄 En curso |
| 5 | Devolver el slice o no | Pendiente |
| 6 | Puntero o valor | Pendiente |
| 7 | El tipo `error` | Pendiente |
| 8 | Paquetes | Pendiente |
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

### Lección 4 — Slice o map (en curso)

*Se completa al aprobar el ejercicio.*
