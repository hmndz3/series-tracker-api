package manejadores

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hmndz3/series-tracker-api/internal/db"
	"github.com/hmndz3/series-tracker-api/internal/modelos"
	"github.com/jackc/pgx/v5"
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

// --------------------------------------------------------------------
// GET /series/{id} — Obtener una serie específica por ID
// --------------------------------------------------------------------

func ObtenerSerie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := obtenerIDDeURL(r)
	if err != nil {
		responderError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	consulta := `
		SELECT id, titulo, descripcion, imagen_url, estado,
		       calificacion, episodios_total, episodios_vistos,
		       creado_en, actualizado_en
		FROM series
		WHERE id = $1
	`

	var s modelos.Serie
	err = db.Pool.QueryRow(ctx, consulta, id).Scan(
		&s.ID, &s.Titulo, &s.Descripcion, &s.ImagenURL, &s.Estado,
		&s.Calificacion, &s.EpisodiosTotal, &s.EpisodiosVistos,
		&s.CreadoEn, &s.ActualizadoEn,
	)

	// Si no hay resultado, la serie no existe → 404
	if errors.Is(err, pgx.ErrNoRows) {
		responderError(w, http.StatusNotFound, "Serie no encontrada")
		return
	}
	if err != nil {
		log.Printf("Error al consultar serie %d: %v", id, err)
		responderError(w, http.StatusInternalServerError, "Error al obtener la serie")
		return
	}

	responderJSON(w, http.StatusOK, s)
}

// --------------------------------------------------------------------
// Helpers compartidos
// --------------------------------------------------------------------

// obtenerIDDeURL extrae el parámetro {id} de la URL y lo convierte a int.
func obtenerIDDeURL(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.Atoi(idStr)
}

// --------------------------------------------------------------------
// EntradaSerie — struct para recibir datos en POST y PUT
// --------------------------------------------------------------------

// EntradaSerie representa el cuerpo (JSON) de las peticiones POST y PUT.
// Todos los campos son punteros para distinguir "no enviado" de "valor vacío".
// Esto es crucial para PUT, donde solo queremos actualizar los campos enviados.
type EntradaSerie struct {
	Titulo          *string `json:"titulo"`
	Descripcion     *string `json:"descripcion"`
	ImagenURL       *string `json:"imagen_url"`
	Estado          *string `json:"estado"`
	Calificacion    *int    `json:"calificacion"`
	EpisodiosTotal  *int    `json:"episodios_total"`
	EpisodiosVistos *int    `json:"episodios_vistos"`
}

// validarCampos verifica las reglas de negocio para los campos enviados.
// Retorna string vacío si todo está bien, o un mensaje de error descriptivo.
func validarCampos(e *EntradaSerie) string {
	if e.Estado != nil && !modelos.EstadosValidos[*e.Estado] {
		return "Estado inválido. Valores permitidos: viendo, completada, pendiente, abandonada"
	}
	if e.Calificacion != nil && (*e.Calificacion < 1 || *e.Calificacion > 10) {
		return "La calificación debe estar entre 1 y 10"
	}
	if e.EpisodiosTotal != nil && *e.EpisodiosTotal < 0 {
		return "El total de episodios no puede ser negativo"
	}
	if e.EpisodiosVistos != nil && *e.EpisodiosVistos < 0 {
		return "Los episodios vistos no pueden ser negativos"
	}
	if e.EpisodiosTotal != nil && e.EpisodiosVistos != nil && *e.EpisodiosVistos > *e.EpisodiosTotal {
		return "Los episodios vistos no pueden superar al total"
	}
	return ""
}

// --------------------------------------------------------------------
// POST /series — Crear una nueva serie
// --------------------------------------------------------------------

func CrearSerie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var entrada EntradaSerie
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido en el cuerpo de la petición")
		return
	}

	// El título es obligatorio para crear una serie
	if entrada.Titulo == nil || strings.TrimSpace(*entrada.Titulo) == "" {
		responderError(w, http.StatusBadRequest, "El título es obligatorio")
		return
	}

	if msgError := validarCampos(&entrada); msgError != "" {
		responderError(w, http.StatusBadRequest, msgError)
		return
	}

	// Valores por defecto si no se envían
	estado := "pendiente"
	if entrada.Estado != nil {
		estado = *entrada.Estado
	}

	episodiosVistos := 0
	if entrada.EpisodiosVistos != nil {
		episodiosVistos = *entrada.EpisodiosVistos
	}

	consulta := `
		INSERT INTO series (titulo, descripcion, imagen_url, estado,
		                    calificacion, episodios_total, episodios_vistos)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, titulo, descripcion, imagen_url, estado,
		          calificacion, episodios_total, episodios_vistos,
		          creado_en, actualizado_en
	`

	var s modelos.Serie
	err := db.Pool.QueryRow(ctx, consulta,
		strings.TrimSpace(*entrada.Titulo),
		entrada.Descripcion,
		entrada.ImagenURL,
		estado,
		entrada.Calificacion,
		entrada.EpisodiosTotal,
		episodiosVistos,
	).Scan(
		&s.ID, &s.Titulo, &s.Descripcion, &s.ImagenURL, &s.Estado,
		&s.Calificacion, &s.EpisodiosTotal, &s.EpisodiosVistos,
		&s.CreadoEn, &s.ActualizadoEn,
	)

	if err != nil {
		log.Printf("Error al crear serie: %v", err)
		responderError(w, http.StatusInternalServerError, "Error al crear la serie")
		return
	}

	// 201 Created es el status correcto para creación exitosa
	responderJSON(w, http.StatusCreated, s)
}

// --------------------------------------------------------------------
// PUT /series/{id} — Actualizar una serie existente
// --------------------------------------------------------------------

func ActualizarSerie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := obtenerIDDeURL(r)
	if err != nil {
		responderError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var entrada EntradaSerie
	if err := json.NewDecoder(r.Body).Decode(&entrada); err != nil {
		responderError(w, http.StatusBadRequest, "JSON inválido en el cuerpo de la petición")
		return
	}

	if msgError := validarCampos(&entrada); msgError != "" {
		responderError(w, http.StatusBadRequest, msgError)
		return
	}

	// COALESCE en SQL: si el valor enviado es NULL, mantiene el valor actual.
	// Así PUT actualiza solo los campos enviados (update parcial).
	consulta := `
		UPDATE series
		SET titulo           = COALESCE($1, titulo),
		    descripcion      = COALESCE($2, descripcion),
		    imagen_url       = COALESCE($3, imagen_url),
		    estado           = COALESCE($4, estado),
		    calificacion     = COALESCE($5, calificacion),
		    episodios_total  = COALESCE($6, episodios_total),
		    episodios_vistos = COALESCE($7, episodios_vistos),
		    actualizado_en   = NOW()
		WHERE id = $8
		RETURNING id, titulo, descripcion, imagen_url, estado,
		          calificacion, episodios_total, episodios_vistos,
		          creado_en, actualizado_en
	`

	var s modelos.Serie
	err = db.Pool.QueryRow(ctx, consulta,
		entrada.Titulo,
		entrada.Descripcion,
		entrada.ImagenURL,
		entrada.Estado,
		entrada.Calificacion,
		entrada.EpisodiosTotal,
		entrada.EpisodiosVistos,
		id,
	).Scan(
		&s.ID, &s.Titulo, &s.Descripcion, &s.ImagenURL, &s.Estado,
		&s.Calificacion, &s.EpisodiosTotal, &s.EpisodiosVistos,
		&s.CreadoEn, &s.ActualizadoEn,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		responderError(w, http.StatusNotFound, "Serie no encontrada")
		return
	}
	if err != nil {
		log.Printf("Error al actualizar serie %d: %v", id, err)
		responderError(w, http.StatusInternalServerError, "Error al actualizar la serie")
		return
	}

	responderJSON(w, http.StatusOK, s)
}

// --------------------------------------------------------------------
// DELETE /series/{id} — Eliminar una serie
// --------------------------------------------------------------------

func EliminarSerie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := obtenerIDDeURL(r)
	if err != nil {
		responderError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	resultado, err := db.Pool.Exec(ctx, `DELETE FROM series WHERE id = $1`, id)
	if err != nil {
		log.Printf("Error al eliminar serie %d: %v", id, err)
		responderError(w, http.StatusInternalServerError, "Error al eliminar la serie")
		return
	}

	if resultado.RowsAffected() == 0 {
		responderError(w, http.StatusNotFound, "Serie no encontrada")
		return
	}

	// 204 No Content: eliminación exitosa, sin cuerpo de respuesta.
	// Es el status estándar REST para DELETE exitoso.
	w.WriteHeader(http.StatusNoContent)
}