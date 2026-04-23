package modelos

import "time"

// Serie representa una serie de TV en el sistema.
// Los campos con `omitempty` se omiten del JSON cuando están vacíos.
type Serie struct {
	ID              int       `json:"id"`
	Titulo          string    `json:"titulo"`
	Descripcion     *string   `json:"descripcion,omitempty"`
	ImagenURL       *string   `json:"imagen_url,omitempty"`
	Estado          string    `json:"estado"`
	Calificacion    *int      `json:"calificacion,omitempty"`
	EpisodiosTotal  *int      `json:"episodios_total,omitempty"`
	EpisodiosVistos int       `json:"episodios_vistos"`
	CreadoEn        time.Time `json:"creado_en"`
	ActualizadoEn   time.Time `json:"actualizado_en"`
}

// EstadosValidos define los únicos valores permitidos para el campo Estado.
var EstadosValidos = map[string]bool{
	"viendo":      true,
	"completada":  true,
	"pendiente":   true,
	"abandonada":  true,
}
