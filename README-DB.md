# Base de Datos Local - Guía de Uso

## Inicio Rápido

### Opción 1: Ejecutar Todo con Docker Compose (Recomendado)

Para ejecutar tanto la base de datos como la aplicación:

```bash
docker-compose up -d
```

Esto iniciará:
- MySQL 8.0 en el puerto 3306
- La aplicación Go en el puerto 8000

La aplicación estará disponible en: `http://localhost:8000`

### Opción 2: Solo Base de Datos

Si solo quieres iniciar MySQL:

```bash
docker-compose up -d db
```

Esto iniciará MySQL 8.0 en el puerto 3306 y ejecutará automáticamente el script `init.sql` para crear todas las tablas necesarias.

### 2. Verificar que MySQL está corriendo

```bash
docker-compose ps
```

Deberías ver el contenedor `accounting_users_db` con estado "Up".

### 3. Conectarte a la base de datos

Puedes conectarte usando cualquier cliente MySQL:

```bash
mysql -h 127.0.0.1 -P 3306 -u preview_usr -plocaldev preview
```

O usando Docker:

```bash
docker exec -it accounting_users_db mysql -u preview_usr -plocaldev preview
```

## Credenciales

- **Host**: localhost (127.0.0.1)
- **Puerto**: 3306
- **Base de datos**: preview
- **Usuario**: preview_usr
- **Contraseña**: localdev
- **Root password**: rootlocal

## Estructura de la Base de Datos

El script `init.sql` crea las siguientes tablas:

1. **users** - Usuarios del sistema
2. **user_code** - Códigos de confirmación/reset password
3. **bookcode** - Códigos de libros
4. **feedback** - Comentarios/feedback de usuarios
5. **DelAccount** - Razones de eliminación de cuentas
6. **metricsterm** - Métricas de términos buscados
7. **sugestedterms** - Términos sugeridos/perdidos
8. **time_used** - Métricas de tiempo usado

## Ejecutar el Script SQL Manualmente

Si necesitas ejecutar el script manualmente (por ejemplo, si el contenedor ya existía antes de agregar el volumen):

```bash
docker exec -i accounting_users_db mysql -u root -prootlocal < init.sql
```

O conectándote directamente:

```bash
mysql -h 127.0.0.1 -P 3306 -u root -prootlocal < init.sql
```

## Reiniciar la Base de Datos

Si necesitas empezar desde cero (⚠️ esto borrará todos los datos):

```bash
docker-compose down -v
docker-compose up -d db
```

## Detener los Servicios

Para detener solo la base de datos:

```bash
docker-compose stop db
```

Para detener todo (app + db):

```bash
docker-compose stop
```

Para detener y eliminar los contenedores (pero mantener los datos):

```bash
docker-compose down
```

## Ejecutar la Aplicación Localmente (sin Docker)

Si prefieres ejecutar la aplicación directamente en tu máquina (no en Docker):

1. Asegúrate de que MySQL esté corriendo:
   ```bash
   docker-compose up -d db
   ```

2. Ejecuta la aplicación:
   ```bash
   go run ./cmd/api/main.go
   ```

La aplicación usará `DATABASE_DSN` de tu archivo `.env` que apunta a `localhost:3306`.

## Notas

- El script `init.sql` se ejecuta automáticamente solo la primera vez que se crea el contenedor.
- Los datos persisten en el volumen Docker `mysql_data`.
- Si ya existe un contenedor con datos, necesitarás ejecutar el script manualmente o eliminar el volumen primero.
