-- Esta tabla almacena todas las series que el usuario va a trackear.

CREATE TABLE IF NOT EXISTS series (
    id               SERIAL PRIMARY KEY,
    titulo           VARCHAR(200) NOT NULL,
    descripcion      TEXT,
    imagen_url       VARCHAR(500),
    estado           VARCHAR(50) DEFAULT 'pendiente',
    calificacion     INTEGER CHECK (calificacion BETWEEN 1 AND 10),
    episodios_total  INTEGER,
    episodios_vistos INTEGER DEFAULT 0,
    creado_en        TIMESTAMP DEFAULT NOW(),
    actualizado_en   TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_series_titulo ON series (LOWER(titulo));

-- Datos ejempl;os para probar 
INSERT INTO series (titulo, descripcion, imagen_url, estado, calificacion, episodios_total, episodios_vistos) VALUES
    ('Breaking Bad', 'Un profesor de química se convierte en fabricante de metanfetaminas.', 'https://image.tmdb.org/t/p/w500/ggFHVNu6YYI5L9pCfOacjizRGt.jpg', 'completada', 10, 62, 62),
    ('The Mandalorian', 'Un cazarrecompensas mandaloriano viaja por la galaxia.', 'https://image.tmdb.org/t/p/w500/sWgBv7LV2PRoQgkxwlibdGXKz1S.jpg', 'viendo', 9, 24, 16),
    ('Chernobyl', 'La historia del desastre nuclear de 1986.', 'https://image.tmdb.org/t/p/w500/hlLXt2tOPT6RRnjiUmoxyG1LTFi.jpg', 'pendiente', NULL, 5, 0);