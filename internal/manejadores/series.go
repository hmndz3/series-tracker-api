package manejadores

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/hmndz3/series-tracker-api/internal/db"
	"github.com/hmndz3/series-tracker-api/internal/modelos"
)

// ListarSeries responde GET /series con la lista de todas las series.
func ListarSeries(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Query SQL: selecciona todas las columnas de todas las series,
	// ordenadas por fecha de creación descendente (las más nuevas primero).
	consulta := `
		SELECT id, titulo, descripcion, imagen_url, estado,
		       calificacion, episodios_total, episodios_vistos,
		       creado_en, actualizado_en
		FROM series
		ORDER BY creado_en DESC
	`

	filas, err := db.Pool.Query(ctx, consulta)
	if err != nil {
		log.Printf("Error al consultar series: %v", err)
		responderError(w, http.StatusInternalServerError, "Error al obtener las series")
		return
	}
	defer filas.Close()

	// Convertir cada fila a un struct Serie
	series := make([]modelos.Serie, 0)
	for filas.Next() {
		var s modelos.Serie
		err := filas.Scan(
			&s.ID,
			&s.Titulo,
			&s.Descripcion,
			&s.ImagenURL,
			&s.Estado,
			&s.Calificacion,
			&s.EpisodiosTotal,
			&s.EpisodiosVistos,
			&s.CreadoEn,
			&s.ActualizadoEn,
		)
		if err != nil {
			log.Printf("Error al leer fila: %v", err)
			responderError(w, http.StatusInternalServerError, "Error al procesar los datos")
			return
		}
		series = append(series, s)
	}

	if err := filas.Err(); err != nil {
		log.Printf("Error al iterar filas: %v", err)
		responderError(w, http.StatusInternalServerError, "Error al procesar los datos")
		return
	}

	responderJSON(w, http.StatusOK, series)
}

// responderJSON es un helper para enviar respuestas JSON con el status correcto.
func responderJSON(w http.ResponseWriter, status int, datos interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(datos); err != nil {
		log.Printf("Error al codificar JSON: %v", err)
	}
}

// responderError envía una respuesta de error en formato JSON estandarizado.
func responderError(w http.ResponseWriter, status int, mensaje string) {
	responderJSON(w, status, map[string]string{
		"error": mensaje,
	})
}