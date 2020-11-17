package graphics

import (
	"fmt"
	"image"
	"image/jpeg"
	"net/http"
	"strings"

	"github.com/bmxguy100/spotify"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"

	log "github.com/sirupsen/logrus"
)

var albumArtCache = ImageCache{
	GetIfInvalid: false,
	GetImage: func(key string) image.Image {
		res, err := http.Get(key)
		if err != nil {
			log.Error(err)
			return nil
		}

		image, err := jpeg.Decode(res.Body)
		if err != nil {
			log.Error(err)
			return nil
		}

		image = resize.Resize(albumArtResolution, albumArtResolution, image, resize.NearestNeighbor)

		return image
	},
}

func getAlbumArt(item *spotify.FullTrack) image.Image {
	imageData := item.Album.Images[0]
	for _, newImageData := range item.Album.Images[1:] {
		if newImageData.Width >= albumArtResolution && newImageData.Width < imageData.Width {
			imageData = newImageData
		}
	}

	return albumArtCache.Retrieve(imageData.URL)
}

const albumArtResolution = 180

const infoAreaPadding = 5
const infoAreaX = albumArtResolution + infoAreaPadding
const infoAreaY = infoAreaPadding
const infoAreaWidth = width - albumArtResolution - 2*infoAreaPadding
const infoAreaHeight = albumArtResolution

const controlAreaPadding = 5
const controlAreaX = controlAreaPadding
const controlAreaY = albumArtResolution + controlAreaPadding
const controlAreaWidth = width - 2*controlAreaPadding
const controlAreaHeight = height - albumArtResolution - 2*controlAreaPadding

func drawSong(context *gg.Context, playerState *spotify.PlayerState) {

	albumArt := getAlbumArt(playerState.Item)
	if albumArt == nil {
		context.SetRGB(.5, .5, .5)
		context.DrawRectangle(0, 0, albumArtResolution, albumArtResolution)
		context.Fill()
	} else {
		context.DrawImage(albumArt, 0, 0)
	}

	var artistNames = make([]string, len(playerState.Item.Artists))
	for i, artist := range playerState.Item.Artists {
		artistNames[i] = artist.Name
	}

	info := fmt.Sprintf("%s\n\u2015\n%s", playerState.Item.Name, strings.Join(artistNames, ", "))

	context.SetRGB(1, 1, 1)
	context.DrawStringWrapped(info, infoAreaX, infoAreaY, 0, 0, infoAreaWidth, 1.5, gg.AlignLeft)

	songPercent := float64(playerState.Progress) / float64(playerState.Item.Duration)
	songPercentX := controlAreaX + controlAreaWidth*songPercent
	songY := controlAreaY + (controlAreaHeight / 2.0)

	context.SetLineWidth(5)

	context.SetRGB(.5, .5, .5)
	context.MoveTo(controlAreaX, songY)
	context.LineTo(controlAreaX+controlAreaWidth, songY)
	context.Stroke()

	context.SetRGB(.25, 1, .25)
	context.MoveTo(controlAreaX, songY)
	context.LineTo(songPercentX, songY)
	context.Stroke()

	context.DrawCircle(songPercentX, songY, 5)
	context.SetRGB(1, 1, 1)
	context.Fill()
}
