# Documentación Swagger - Accounting Users API

## Instalación

### 1. Instalar Swagger CLI

**Fuera del contenedor Docker:**
```bash
go install github.com/swaggo/swag/cmd/swag@v1.8.1
```

O usando el Makefile:
```bash
make swagger-install
```

**Dentro del contenedor Docker:**
```bash
go install github.com/swaggo/swag/cmd/swag@v1.8.1
```

### 2. Instalar dependencias del proyecto

```bash
go mod tidy
```

## Generar Documentación

### Opción 1: Usando Make (recomendado fuera del contenedor)

```bash
make swagger
```

### Opción 2: Comando directo (si no tienes make o estás en Docker)

```bash
swag init -g cmd/api/main.go -o docs
```

**Si estás dentro de un contenedor Docker y no tienes `swag` instalado:**

```bash
# Instalar swag primero (usar versión compatible con go.mod)
go install github.com/swaggo/swag/cmd/swag@v1.8.1

# Luego generar docs
swag init -g cmd/api/main.go -o docs

# O en una sola línea:
go install github.com/swaggo/swag/cmd/swag@v1.8.1 && swag init -g cmd/api/main.go -o docs
```

## Ejecutar la API

```bash
make run
# o
go run ./cmd/api/main.go
```

## Acceder a Swagger UI

Una vez que la API esté corriendo, accede a:

- **Swagger UI**: http://localhost:8000/swagger/index.html
- **JSON de Swagger**: http://localhost:8000/swagger/doc.json

## Agregar Documentación a Nuevos Endpoints

Para documentar un nuevo endpoint, agrega comentarios Swagger antes de la función handler:

```go
// nombreFuncion godoc
// @Summary      Descripción breve
// @Description  Descripción detallada
// @Tags         tag-name
// @Accept       json
// @Produce      json
// @Param        param  path/query/body  type  required  "Descripción"
// @Success      200    {object}         map[string]interface{}
// @Failure      400    {object}         map[string]interface{}
// @Router       /ruta/{param} [method]
func (app *application) nombreFuncion(w http.ResponseWriter, r *http.Request) {
    // código del handler
}
```

### Ejemplo Completo

```go
// getUser godoc
// @Summary      Obtener usuario por ID
// @Description  Retorna la información de un usuario específico por su ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del usuario"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Router       /v1/user/{id} [get]
func (app *application) getUser(w http.ResponseWriter, r *http.Request) {
    // código...
}
```

## Tags Disponibles

Los endpoints están organizados en los siguientes tags:

- `status` - Estado de la API
- `users` - Gestión de usuarios
- `auth` - Autenticación
- `metrics` - Métricas y estadísticas
- `comments` - Comentarios y feedback
- `books` - Códigos de libros

## Regenerar Documentación

Cada vez que agregues o modifiques comentarios Swagger, debes regenerar la documentación:

```bash
make swagger
# o directamente:
swag init -g cmd/api/main.go -o docs
```

Luego reinicia la aplicación para ver los cambios.

## Solución de Problemas

### Error: "make: not found" (dentro de Docker)

Ejecuta el comando directamente:
```bash
swag init -g cmd/api/main.go -o docs
```

Si no tienes `swag` instalado:
```bash
go install github.com/swaggo/swag/cmd/swag@v1.8.1 && swag init -g cmd/api/main.go -o docs
```

### Error: "swag: command not found"

Asegúrate de que `$GOPATH/bin` o `$HOME/go/bin` esté en tu `PATH`:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
# o
export PATH=$PATH:$HOME/go/bin
```

## Notas

- Los archivos generados en `docs/` están en `.gitignore` y deben regenerarse en cada máquina
- El endpoint `/swagger/*` está configurado para servir la documentación
- La documentación se genera automáticamente desde los comentarios en el código
- Si estás usando Docker, considera generar la documentación antes de construir la imagen, o montar el directorio `docs/` como volumen
