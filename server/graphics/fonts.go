package graphics

import (
	"io/ioutil"

	"github.com/AndreKR/multiface"
	"github.com/golang/freetype/truetype"
	"github.com/markbates/pkger"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	log "github.com/sirupsen/logrus"
)

func loadFonts() (face *multiface.Face, err error) {
	face = new(multiface.Face)

	err = loadTTFFont(face, pkger.Include("/fonts/NotoSans-Regular.ttf"))
	if err != nil {
		return
	}

	err = loadOTFFont(face, pkger.Include("/fonts/NotoSansCJKkr-Regular.otf"))
	if err != nil {
		return
	}

	return
}

func loadTTFFont(mface *multiface.Face, path string) (err error) {
	file, err := pkger.Open(path)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	ttffont, err := truetype.Parse(data)
	if err != nil {
		return
	}

	face := truetype.NewFace(ttffont, &truetype.Options{Size: 15})
	mface.AddTruetypeFace(face, ttffont)
	log.Printf("Loaded %s\n", path)
	return
}

func loadOTFFont(mface *multiface.Face, path string) (err error) {
	file, err := pkger.Open(path)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	otffont, err := opentype.Parse(data)
	if err != nil {
		return
	}

	face, err := opentype.NewFace(otffont, &opentype.FaceOptions{
		Size: 15,
		DPI: 72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return
	}

	mface.AddFace(face)
	log.Printf("Loaded %s\n", path)
	return
}
