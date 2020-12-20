package main

import (
	"os"

	"github.com/bmxguy100/spotify-screen/api"
	"github.com/bmxguy100/spotify-screen/graphics"
	"github.com/bmxguy100/spotify-screen/serial"

	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
	}

	err := serial.InitSerial()
	if err != nil {
		log.WithError(err).Fatal("Error connecting to arduino")
	}

	go api.SpotifyServer()

	graphics.FrameGenerator()
}
