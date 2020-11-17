package main

import (
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

	err := serial.InitSerial()
	if err != nil {
		log.Fatal(err)
	}
	
	go api.SpotifyServer()
	
	graphics.FrameGenerator()
}
