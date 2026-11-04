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
| 3 | Tipo con nombre: `type Priority int` | 🔄 En curso |
| 4 | Slice o map | Pendiente |
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

### Lección 3 — Tipo con nombre (en curso)

*Se completa al aprobar el ejercicio.*
