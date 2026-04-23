package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/hmndz3/series-tracker-api/internal/db"
)

func main() {
	// Conectar a PostgreSQL antes de arrancar el servidor.
	// Si falla la conexión, no tiene sentido seguir.
	if err := db.Conectar(); err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer db.Cerrar()

	// Puerto: Railway inyecta PORT como variable de entorno.
	// En local, si no existe, usamos 8080.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := chi.NewRouter()

	// Middlewares básicos
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS: durante desarrollo permitimos todo.
	// Para producción más segura, reemplazar "*" con la URL del frontend.
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Endpoint raíz: confirma que el servidor está corriendo
	r.Get("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Series Tracker API",
			"status":  "ok",
		})
	})

	// Endpoint de salud: verifica que la BD responde
	r.Get("/salud", func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		if err := db.Pool.Ping(ctx); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			json.NewEncoder(w).Encode(map[string]string{
				"status": "error",
				"error":  "base de datos no disponible",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"status":        "ok",
			"base_de_datos": "conectada",
		})
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Servidor arrancando en %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Servidor falló: %v", err)
	}
}