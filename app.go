package main

import (
	"context"
	"database/sql"
	"log"
	"os/exec"
	"patio-desktop/database"
	"patio-desktop/models"
	"patio-desktop/services"
)

// App struct
type App struct {
	ctx       context.Context
	db        *sql.DB
	go2rtcCmd *exec.Cmd
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	db, err := database.InitDB("./patio.db")
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}
	a.db = db
	cameras, _ := database.GetCameras(db)
	services.GenerateGo2rtcConfig(cameras)

	a.go2rtcCmd = exec.Command("./bin/go2rtc/go2rtc")
	err = a.go2rtcCmd.Start()
	if err != nil {
		log.Println("Error starting go2rtc: ", err)
	}
}

func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
	if a.go2rtcCmd != nil && a.go2rtcCmd.Process != nil {
		a.go2rtcCmd.Process.Kill()
	}
}

// GetCameras retourne la liste des caméras (appelable depuis Svelte)
func (a *App) GetCameras() ([]database.Camera, error) {
	return database.GetCameras(a.db)
}

func (a *App) AddCamera(name, rtspURL string) (int64, error) {
	id, err := database.AddCamera(a.db, name, rtspURL)
	if err != nil {
		return 0, err
	}

	err = services.AddCameraToGo2rtc(id, rtspURL)
	if err != nil {
		log.Println("Error adding camera to go2rtc: ", err)
	}

	return id, nil
}

// GetZonesByCamera retourne les zones d'une caméra (appelable depuis Svelte)
func (a *App) GetZonesByCamera(cameraID int64) ([]database.Zone, error) {
	return database.GetZonesByCamera(a.db, cameraID)
}

// AddZone ajoute une zone à une caméra (appelable depuis Svelte)
func (a *App) AddZone(z database.Zone) (int64, error) {
	return database.AddZone(a.db, z)
}

// GetEventsByZone retourne l'historique des events d'une zone (appelable depuis Svelte)
func (a *App) GetEventsByZone(zoneID int64) ([]database.Event, error) {
	return database.GetEventsByZone(a.db, zoneID)
}

// ScanNetwork discovers ONVIF cameras on the network (appelable depuis Svelte)
func (a *App) ScanNetwork() ([]models.DiscoveredCamera, error) {
	return services.ScanNetwork()
}

func (a *App) GetRTSPStreamURL(endpoint, username, password string) (string, error) {
	return services.GetRTSPStreamURL(endpoint, username, password)
}

func (a *App) GenerateGo2rtcConfig(listeCamera []database.Camera) string {
	return services.GenerateGo2rtcConfig(listeCamera)
}

// AddCameraToGo2rtc enregistre le flux d'une caméra auprès de go2rtc (appelable depuis Svelte)
func (a *App) AddCameraToGo2rtc(camID int64, rtspURL string) error {
	return services.AddCameraToGo2rtc(camID, rtspURL)
}
