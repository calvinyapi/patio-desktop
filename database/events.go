package database

import (
	"database/sql"
	"strconv"
)

type Event struct {
	ID              int64  `json:"id"`
	ZoneID          int64  `json:"zone_id"`
	EventType       string `json:"event_type"`
	PersonCount     int    `json:"person_count"`
	DurationSeconds int    `json:"duration_seconds"`
	Timestamp       string `json:"timestamp"`
	Synced          bool   `json:"synced"`
}

// AddEvent insère un nouvel event, détecté par l'inférence (YOLO), pour une zone donnée.
func AddEvent(db *sql.DB, e Event) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO events (zone_id, event_type, person_count, duration_seconds)
		 VALUES (?, ?, ?, ?)`,
		e.ZoneID, e.EventType, e.PersonCount, e.DurationSeconds,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetUnsyncedEvents retourne tous les events pas encore envoyés à PATIO-API.
// C'est cette fonction que le ticker de sync appellera périodiquement.
func GetUnsyncedEvents(db *sql.DB) ([]Event, error) {
	rows, err := db.Query(
		`SELECT id, zone_id, event_type, person_count, duration_seconds, timestamp, synced
		 FROM events WHERE synced = FALSE ORDER BY timestamp`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.ZoneID, &e.EventType, &e.PersonCount, &e.DurationSeconds, &e.Timestamp, &e.Synced); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// GetEventsByZone retourne l'historique des events pour une zone donnée (utile pour l'affichage local).
func GetEventsByZone(db *sql.DB, zoneID int64) ([]Event, error) {
	rows, err := db.Query(
		`SELECT id, zone_id, event_type, person_count, duration_seconds, timestamp, synced
		 FROM events WHERE zone_id = ? ORDER BY timestamp DESC`, zoneID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err := rows.Scan(&e.ID, &e.ZoneID, &e.EventType, &e.PersonCount, &e.DurationSeconds, &e.Timestamp, &e.Synced); err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, rows.Err()
}

// MarkEventsSynced marque une liste d'events comme envoyés avec succès à PATIO-API.
func MarkEventsSynced(db *sql.DB, ids []int64) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`UPDATE events SET synced = TRUE WHERE id = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, id := range ids {
		if _, err := stmt.Exec(id); err != nil {
			return err
		}
	}

	return tx.Commit()
}

// DeleteSyncedEventsOlderThan fait le ménage des events déjà envoyés et anciens
// (garde SQLite comme un buffer léger, pas un entrepôt définitif).
func DeleteSyncedEventsOlderThan(db *sql.DB, days int) error {
	_, err := db.Exec(
		`DELETE FROM events WHERE synced = TRUE AND timestamp < datetime('now', ?)`,
		"-"+strconv.Itoa(days)+" days",
	)
	return err
}
