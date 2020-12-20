package graphics

import (
	"image"
	"time"

	"github.com/bmxguy100/spotify-screen/api"
	"github.com/fogleman/gg"
	"github.com/markbates/pkger"

	log "github.com/sirupsen/logrus"
)

const realWidth = 320
const width = 310
const height = 240
const minFrameTime = time.Millisecond * 100

var FrameChannel = make(chan image.Image)

func FrameGenerator() {
	var err error
	spotifyLogoFile, err := pkger.Open("/img/spotify_icon.png")
	if err != nil {
		log.WithError(err).Fatal("Error loading '/img/spotify_icon.png'")
	}
	spotifyLogo, _, err = image.Decode(spotifyLogoFile)
	if err != nil {
		log.WithError(err).Fatal("Error decoding '/img/spotify_icon.png'")
	}

	face, err := loadFonts()
	if err != nil {
		log.WithError(err).Fatal("Error loading fonts")
	}

	nextFrame := time.Now()
	for {
		time.Sleep(time.Until(nextFrame))
		startTime := time.Now()
		nextFrame = time.Now().Add(minFrameTime)

		context := gg.NewContext(realWidth, height)
		context.SetFontFace(face)

		context.SetRGB(.3, .3, .3)
		context.Clear()

		state := <-api.PlaybackStateChannel

		if state.Err != nil {
			log.WithError(err).Error("Error in API")
		} else if !state.IsAuthenticated {
			log.Debug("Displaying Unauthenticated")
			drawUnauthenticated(context, state.AuthUrl)
		} else if state.State.CurrentlyPlayingType == "ad" {
			log.Debug("Displaying Ad")
			drawAd(context)
		} else if state.State.Item != nil {
			log.Debug("Displaying Song")
			drawSong(context, &state.State)
		} else {
			log.Debug("Displaying Nothing")
			drawNothing(context)
		}
		log.WithField("Time", time.Since(startTime)).Debug("Drew frame")

		FrameChannel <- context.Image()
	}
}

func drawUnauthenticated(context *gg.Context, url string) {
	context.SetRGB(1, 1, 1)
	context.DrawStringAnchored(url, width/2, height/2, 0.5, 0.5)
}

var spotifyLogo image.Image

func drawNothing(context *gg.Context) {
	context.DrawImageAnchored(spotifyLogo, width/2, height/2, 0.5, 0.5)
}

func drawAd(context *gg.Context) {
	context.SetRGB(1, 1, 1)
	context.DrawStringAnchored("ad", width/2, height/2, 0.5, 0.5)
}
