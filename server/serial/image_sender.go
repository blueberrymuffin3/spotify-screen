package serial

import (
	"image"
	"os"
	"time"

	"github.com/tarm/serial"

	log "github.com/sirupsen/logrus"
)

const imageWidth = 240
const imageHeight = 320
const blockSize = 240
const blockCount = imageWidth * imageHeight / (blockSize / 2)

var reqChannel = make(chan bool)
var serialPort *serial.Port

var lastImage image.Image

func SendFrame(image image.Image) (err error) {
	startTime := time.Now()

	isSame := compareImages(image, lastImage)

	lastImage = image

	if isSame {
		log.Info("Not sending frame, no changes")
		return
	}

	buffer := encodeAndRotateImage(image)

	header := make([]byte, 8)
	x1 := uint16(0)
	y1 := uint16(0)
	x2 := uint16(240)
	y2 := uint16(320)

	header[0] = uint8(x1)
	header[1] = uint8(x1 >> 8)
	header[2] = uint8(y1)
	header[3] = uint8(y1 >> 8)
	header[4] = uint8(x2)
	header[5] = uint8(x2 >> 8)
	header[6] = uint8(y2)
	header[7] = uint8(y2 >> 8)

	_, err = serialPort.Write(header)
	if err != nil {
		lastImage = nil
		return
	}
	_, err = serialPort.Write(buffer)
	if err != nil {
		lastImage = nil
		return
	}
	log.WithField("Time", time.Since(startTime)).WithField("Bytes", len(buffer)).Info("Send frame")
	return
}

func compareImages(imageA image.Image, imageB image.Image) bool {
	if imageA == nil || imageB == nil {
		return false
	}

	return false

	return true
}

func encodeAndRotateImage(source image.Image) []byte {
	sourceSize := source.Bounds().Size()
	buff := make([]byte, sourceSize.X*sourceSize.Y*2)

	for y := 0; y < sourceSize.Y; y++ {
		for x := 0; x < sourceSize.X; x++ {
			i := (x*sourceSize.Y + (sourceSize.Y - 1 - y)) * 2
			sr, sg, sb, _ := source.At(x, y).RGBA()
			r := (byte)(sr >> 11)
			g := (byte)(sg >> 10)
			b := (byte)(sb >> 11)
			buff[i] = (byte)((g << 5) | b)
			buff[i+1] = (byte)((r << 3) | (g >> 3))
		}
	}

	return buff
}

func echo() {
	for {
		buff := make([]byte, 256)
		len, _ := serialPort.Read(buff)
		if len == 0 {
			continue
		}

		os.Stdout.Write(buff[0:len])
	}
}

func InitSerial() (err error) {
	serialPort, err = serial.OpenPort(&serial.Config{
		Name:        "/dev/serial/by-id/usb-Arduino__www.arduino.cc__0043_64934333235351700140-if00",
		Baud:        4000000,
		ReadTimeout: time.Millisecond * 250,
	})

	if err != nil {
		return
	}

	err = serialPort.Flush()

	go echo()
	return
}
