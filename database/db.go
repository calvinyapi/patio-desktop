package database

import (
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

// InitDB ouvre (ou crée) le fichier SQLite et s'assure que le schéma existe.
func InitDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("ouverture de la base: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connexion à la base: %w", err)
	}

	if err := createSchema(db); err != nil {
		return nil, fmt.Errorf("création du schéma: %w", err)
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	schema := `
	PRAGMA foreign_keys = ON;

	CREATE TABLE IF NOT EXISTS cameras (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		rtsp_url TEXT,
		is_active BOOLEAN DEFAULT true,
		status TEXT DEFAULT 'offline',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS zones (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		camera_id INTEGER NOT NULL,
		name TEXT,
		x1 INTEGER NOT NULL,
		y1 INTEGER NOT NULL,
		x2 INTEGER NOT NULL,
		y2 INTEGER NOT NULL,
		threshold_seconds INTEGER DEFAULT 5,
		is_critical BOOLEAN DEFAULT true,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (camera_id) REFERENCES cameras(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		zone_id INTEGER NOT NULL,
		event_type TEXT NOT NULL,
		person_count INTEGER,
		duration_seconds INTEGER,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
		synced BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (zone_id) REFERENCES zones(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_events_synced ON events(synced);
	CREATE INDEX IF NOT EXISTS idx_zones_camera ON zones(camera_id);
	`

	if _, err := db.Exec(schema); err != nil {
		return err
	}

	// Migrations best-effort : CREATE TABLE IF NOT EXISTS ne touche pas une
	// table déjà créée. On ajoute les colonnes manquantes et on ignore
	// l'erreur "duplicate column name" quand elles existent déjà.
	migrations := []string{
		`ALTER TABLE zones ADD COLUMN is_critical BOOLEAN DEFAULT true`,
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil && !strings.Contains(err.Error(), "duplicate column name") {
			return err
		}
	}
	return nil
}
