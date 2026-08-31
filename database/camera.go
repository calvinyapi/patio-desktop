package database

import (
	"database/sql"
)

type Camera struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	RtspURL  string `json:"rtsp_url"`
	IsActive bool   `json:"is_active"`
	Status   string `json:"status"`
}

// AddCamera insère une nouvelle caméra et retourne son ID généré.
func AddCamera(db *sql.DB, name, rtspURL string) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO cameras (name, rtsp_url) VALUES (?, ?)`,
		name, rtspURL,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetCameras retourne toutes les caméras enregistrées.
func GetCameras(db *sql.DB) ([]Camera, error) {
	rows, err := db.Query(
		`SELECT id, name, rtsp_url, is_active, status FROM cameras ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cameras []Camera
	for rows.Next() {
		var c Camera
		if err := rows.Scan(&c.ID, &c.Name, &c.RtspURL, &c.IsActive, &c.Status); err != nil {
			return nil, err
		}
		cameras = append(cameras, c)
	}
	return cameras, rows.Err()
}

// GetCameraByID retourne une seule caméra par son ID.
func GetCameraByID(db *sql.DB, id int64) (*Camera, error) {
	var c Camera
	err := db.QueryRow(
		`SELECT id, name, rtsp_url, is_active, status FROM cameras WHERE id = ?`, id,
	).Scan(&c.ID, &c.Name, &c.RtspURL, &c.IsActive, &c.Status)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// UpdateCameraStatus met à jour uniquement le statut d'une caméra (online/offline).
func UpdateCameraStatus(db *sql.DB, id int64, status string) error {
	_, err := db.Exec(
		`UPDATE cameras SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		status, id,
	)
	return err
}

// UpdateCamera met à jour le nom et l'URL RTSP d'une caméra.
func UpdateCamera(db *sql.DB, id int64, name, rtspURL string) error {
	_, err := db.Exec(
		`UPDATE cameras SET name = ?, rtsp_url = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		name, rtspURL, id,
	)
	return err
}

// DeleteCamera supprime une caméra (et ses zones/events via ON DELETE CASCADE).
func DeleteCamera(db *sql.DB, id int64) error {
	_, err := db.Exec(`DELETE FROM cameras WHERE id = ?`, id)
	return err
}
