package services

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"patio-desktop/database"
)

func GenerateGo2rtcConfig(listeCamera []database.Camera) string {
	f, err := os.Create("./go2rtc.yaml")
	if err != nil {
		fmt.Println(err)
	}
	f.WriteString("streams:\n")
	for _, cam := range listeCamera {
		fmt.Println("La caméra", cam.Name, "a été ajouté")
		line := fmt.Sprintf("  camera_%d: %s\n", cam.ID, cam.RtspURL)
		f.WriteString(line)
	}
	fmt.Println("Le café c'est bon")
	f.Close()
	return "La fonction semble s'être bien passé"
}

func AddCameraToGo2rtc(camID int64, rtspURL string) error {
	url := fmt.Sprintf("http://localhost:1984/api/streams?src=%s&name=camera_%d",
		url.QueryEscape(rtspURL), camID)

	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
