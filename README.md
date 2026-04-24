# Series Tracker — BACKEND
-----------------------------------------------------------------------------------------------------------------------------------------------------
Backend REST para el proyecto Series Tracker. Expone una API que permite gestionar una colección personal de series.

**Frontend desplegado:** https://series-tracker-teal.vercel.app  
**API desplegada:** https://series-tracker-api-production.up.railway.app  
**Repo del frontend:** https://github.com/hmndz3/series-tracker-client

-----------------------------------------------------------------------------------------------------------------------------------------------------

## Screenshots

![Página principal](docs/screenshots/paginaprincipal.png)

![Detalle de una serie](docs/screenshots/fichaserie.png)

![Formulario de nueva serie](docs/screenshots/nuevaserie.png)

------------------------------------------------------------------------------------------------------------------------------------------------------

## Stack

- Go 1.25
- Router HTTP: chi/v5
- Base de datos: PostgreSQL 16
- Driver DB: pgx/v5
- Hosting: Railway

-----------------------------------------------------------------------------------------------------------------------------------------------------

## Endpoints

GET    - `/`               - Estado del servidor                  
GET    - `/salud`          - Verifica conexión con la BD          
GET    - `/series`         - Listar series (con filtros opcionales) 
GET    - `/series/{id}`    - Obtener una serie por ID             
POST   - `/series`         - Crear una serie nueva                
PUT    - `/series/{id}`    - Actualizar una serie existente       
DELETE - `/series/{id}`    - Eliminar una serie                   

### Query parameters de `GET /series`

`q`       - Búsqueda por título (case-insensitive)              
`sort`    - Campo de ordenamiento: `titulo`, `calificacion`, `creado_en` - `creado_en`  
`order`   - Dirección: `asc` o `desc`                            - `desc`       
`page`    - Número de página                                     - `1`          
`limit`   - Resultados por página (máximo 100)                   - `10`         

-----------------------------------------------------------------------------------------------------------------------------------------------------

## Sobre CORS

Como el cliente (Vercel) y el servidor (Railway) corren en dominios distintos, el navegador bloquearía las peticiones `fetch()` por seguridad. Para permitirlo, el servidor envía estos headers:
-Access-Control-Allow-Origin: *
-Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
-Access-Control-Allow-Headers: Content-Type

Durante desarrollo se permite cualquier origen (`*`).

-----------------------------------------------------------------------------------------------------------------------------------------------------

## Correr localmente (Windows)

### Requisitos

- Go 1.25 o superior
- PostgreSQL corriendo localmente o URL de una instancia remota

### Pasos

1. Clonar el repo:
```powershell
   git clone https://github.com/hmndz3/series-tracker-api.git
   cd series-tracker-api
```

2. Instalar dependencias:
```powershell
   go mod download
```

3. Configurar la variable de entorno con la URL de PostgreSQL:
```powershell
   $env:DATABASE_URL="postgresql://usuario:password@localhost:5432/series_tracker"
```

4. Ejecutar la migración inicial (crea la tabla y datos de ejemplo). Abrir `migrations/001_crear_tabla_series.sql` y ejecutarlo en tu cliente de PostgreSQL (pgAdmin, DBeaver, etc.).

5. Arrancar el servidor:
```powershell
   go run main.go
```

6. El servidor queda disponible en `http://localhost:8080`.

-----------------------------------------------------------------------------------------------------------------------------------------------------

## Estructura del proyecto

series-tracker-api/
├── main.go
├── go.mod
├── internal/
│   ├── db/
│   │   └── db.go
│   ├── modelos/
│   │   └── serie.go
│   └── manejadores/
│       └── series.go
├── migrations/
│   └── 001_crear_tabla_series.sql
└── docs/
└── screenshots/

-----------------------------------------------------------------------------------------------------------------------------------------------------

## Challenges implementados

- Códigos HTTP correctos (+20) — 200, 201, 204, 400, 404
- Validación server-side con respuestas JSON descriptivas (+20)
- Paginación con `?page=` y `?limit=` (+30)
- Búsqueda por nombre con `?q=` (+15)
- Ordenamiento con `?sort=` y `?order=asc|desc` (+15)

Total backend: 100 puntos de challenges técnicos.

-----------------------------------------------------------------------------------------------------------------------------------------------------

## Reflexión

Este proyecto me gustó más de lo que esperaba. No fue difícil, pero sí fue raro porque tenía que pensar en dos cosas al mismo tiempo: el backend como un servicio aparte y el cliente como algo totalmente separado que solo consumía datos. En el Lab 5 todo estaba mezclado, así que ahora tocó acostumbrarme a pensar en términos de "contratos" entre ambos lados.

Con Go ya había trabajado antes pero muy poco, así que al principio tuve que acordarme de cosas como la sintaxis y cómo manejar errores. Una vez que agarré el ritmo, me sentí cómodo y me di cuenta que para APIs es bastante directo. Los paquetes `chi` y `pgx` me sirvieron bastante y no tuve que pelearme mucho con ellos.

La parte que más me entretuvo fue organizar el código en carpetas separadas (`db`, `modelos`, `manejadores`). Al principio tuve un error medio tonto porque puse archivos en el paquete equivocado, pero después de arreglarlo todo quedó bastante ordenado y me fue fácil encontrar las cosas.

Lo que más me sorprendió fue el deploy en Railway. Pensé que iba a ser complicado pero conecté el repo y prácticamente se encargó de todo solo. También fue útil que detectara automáticamente el proyecto y me conectara la base de datos sin tener que copiar credenciales a mano.

Sí volvería a usar este stack para otro proyecto parecido. Me quedó la curiosidad de probar algo como Gin en el futuro, porque al final sí hay cosas repetitivas que seguro ya están resueltas en frameworks más grandes.
-----------------------------------------------------------------------------------------------------------------------------------------------------