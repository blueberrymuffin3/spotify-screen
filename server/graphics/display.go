package graphics

import (
	"github.com/bmxguy100/spotify-screen/api"
	"github.com/bmxguy100/spotify-screen/serial"
	"github.com/fogleman/gg"

	log "github.com/sirupsen/logrus"
)

const realWidth = 320
const width = 310
const height = 240

func FrameGenerator() {
	context := gg.NewContext(realWidth, height)

	face, err := loadFonts()
	if err != nil {
		log.Fatal(err)
	}

	context.SetFontFace(face)

	for {
		context.SetRGB(.3, .3, .3)
		context.Clear()

		state := <-api.PlaybackStateChannel

		context.Push()
		if state.Err != nil {
			log.Error(err)
		} else if !state.IsAuthenticated {
			log.Info("Displaying Unauthenticated")
			drawUnauthenticated(context, state.AuthUrl)
		} else if state.State.CurrentlyPlayingType == "ad" {
			log.Info("Displaying Ad")
			drawAd(context)
		} else if state.State.Item != nil {
			log.Info("Displaying Song")
			drawSong(context, &state.State)
		} else {
			log.Info("Displaying Nothing")
			drawNothing(context)
		}
		context.Pop()

		err = serial.SendFrame(context.Image())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func drawUnauthenticated(context *gg.Context, url string) {
	context.SetRGB(1, 1, 1)
	context.DrawStringAnchored(url, width/2, height/2, 0.5, 0.5)
}

func drawNothing(context *gg.Context) {

}

func drawAd(context *gg.Context) {
	context.SetRGB(1, 1, 1)
	context.DrawStringAnchored("ad", width/2, height/2, 0.5, 0.5)
}
