package database

import (
	"database/sql"
)

type Zone struct {
	ID               int64  `json:"id"`
	CameraID         int64  `json:"camera_id"`
	Name             string `json:"name"`
	X1               int    `json:"x1"`
	Y1               int    `json:"y1"`
	X2               int    `json:"x2"`
	Y2               int    `json:"y2"`
	ThresholdSeconds int    `json:"threshold_seconds"`
}

// AddZone insère une nouvelle zone pour une caméra donnée.
func AddZone(db *sql.DB, z Zone) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO zones (camera_id, name, x1, y1, x2, y2, threshold_seconds)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		z.CameraID, z.Name, z.X1, z.Y1, z.X2, z.Y2, z.ThresholdSeconds,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetZonesByCamera retourne toutes les zones d'une caméra donnée.
func GetZonesByCamera(db *sql.DB, cameraID int64) ([]Zone, error) {
	rows, err := db.Query(
		`SELECT id, camera_id, name, x1, y1, x2, y2, threshold_seconds
		 FROM zones WHERE camera_id = ? ORDER BY id`, cameraID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var zones []Zone
	for rows.Next() {
		var z Zone
		if err := rows.Scan(&z.ID, &z.CameraID, &z.Name, &z.X1, &z.Y1, &z.X2, &z.Y2, &z.ThresholdSeconds); err != nil {
			return nil, err
		}
		zones = append(zones, z)
	}
	return zones, rows.Err()
}

// UpdateZone met à jour les coordonnées et le seuil d'une zone.
func UpdateZone(db *sql.DB, z Zone) error {
	_, err := db.Exec(
		`UPDATE zones SET name = ?, x1 = ?, y1 = ?, x2 = ?, y2 = ?, threshold_seconds = ?
		 WHERE id = ?`,
		z.Name, z.X1, z.Y1, z.X2, z.Y2, z.ThresholdSeconds, z.ID,
	)
	return err
}

// DeleteZone supprime une zone (et ses events via ON DELETE CASCADE).
func DeleteZone(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM zones WHERE id = ?`, id)
	return err
}
