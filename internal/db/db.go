package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pool es el pool global de conexiones a PostgreSQL.
// Se inicializa una sola vez al arrancar el servidor.
var Pool *pgxpool.Pool

// Conectar establece la conexión inicial con la base de datos.
func Conectar() error {
	urlBaseDatos := os.Getenv("DATABASE_URL")
	if urlBaseDatos == "" {
		return fmt.Errorf("DATABASE_URL no está definida en las variables de entorno")
	}

	// Configuración del pool
	config, err := pgxpool.ParseConfig(urlBaseDatos)
	if err != nil {
		return fmt.Errorf("error al parsear DATABASE_URL: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// Crear el pool
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("error al crear pool de conexiones: %w", err)
	}

	// Verificar que la conexión funciona con un ping
	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("error al hacer ping a la base de datos: %w", err)
	}

	Pool = pool
	log.Println("Conexión a PostgreSQL establecida correctamente")
	return nil
}

// Cerrar libera todas las conexiones del pool.
func Cerrar() {
	if Pool != nil {
		Pool.Close()
		log.Println("Pool de conexiones cerrado")
	}
}