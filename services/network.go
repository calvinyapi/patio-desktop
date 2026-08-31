package services

import (
	"context"
	"fmt"
	"net/url"
	"patio-desktop/models"
	"time"

	"github.com/0x524a/onvif-go"
	"github.com/0x524a/onvif-go/discovery"
)

func ScanNetwork() ([]models.DiscoveredCamera, error) {
	ctx := context.Background()
	devices, err := discovery.Discover(ctx, 5*time.Second)

	if err != nil {
		fmt.Printf("Error discovering ONVIF devices: %v\n", err)
		return nil, err
	}

	var result []models.DiscoveredCamera

	for _, d := range devices {
		name, err := url.QueryUnescape(d.GetName())
		if err != nil {
			fmt.Printf("Error unescaping camera name: %v\n", err)
			name = d.GetName()
		}
		result = append(result, models.DiscoveredCamera{
			Name:     name,
			Endpoint: d.GetDeviceEndpoint(),
		})
	}

	return result, nil
}

func GetRTSPStreamURL(endpoint, username, password string) (string, error) {
	ctx := context.Background()

	client, err := onvif.NewClient(
		endpoint,
		onvif.WithCredentials(username, password),
		onvif.WithTimeout(15*time.Second),
	)
	if err != nil {
		return "", fmt.Errorf("connexion caméra: %w", err)
	}

	profiles, err := client.GetProfiles(ctx)
	if err != nil {
		return "", fmt.Errorf("récupération profils: %w", err)
	}
	if len(profiles) == 0 {
		return "", fmt.Errorf("aucun profil média trouvé sur cette caméra")
	}

	streamURI, err := client.GetStreamURI(ctx, profiles[0].Token)
	if err != nil {
		return "", fmt.Errorf("récupération flux RTSP: %w", err)
	}

	return streamURI.URI, nil
}
