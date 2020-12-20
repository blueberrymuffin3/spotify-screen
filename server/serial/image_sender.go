package serial

import (
	"image"
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

var lastFrame image.Image

func SendFrame(frame image.Image) (err error) {
	startTime := time.Now()

	shouldUpdate, updateRegion := getUpdateRegion(lastFrame, frame)

	lastFrame = frame

	if !shouldUpdate {
		log.Debug("Not sending frame, no changes")
		return
	}

	buffer := encodeAndRotateImage(frame, updateRegion)

	header := make([]byte, 8)

	x1 := uint16(updateRegion.Min.X)
	y1 := uint16(updateRegion.Min.Y)
	x2 := uint16(updateRegion.Max.X)
	y2 := uint16(updateRegion.Max.Y)

	header[0] = uint8(y1)
	header[1] = uint8(y1 >> 8)
	header[2] = uint8(x1)
	header[3] = uint8(x1 >> 8)
	header[4] = uint8(y2)
	header[5] = uint8(y2 >> 8)
	header[6] = uint8(x2)
	header[7] = uint8(x2 >> 8)

	_, err = serialPort.Write(header)
	if err != nil {
		lastFrame = nil
		return
	}
	_, err = serialPort.Write(buffer)
	if err != nil {
		lastFrame = nil
		return
	}

	isOk, err := confirmSend()
	if err != nil {
		log.WithError(err).Error("Error confirming send")
	}
	if err != nil || !isOk {
		lastFrame = nil
		serialPort.Flush()
		time.Sleep(3 * time.Second)
	}

	log.WithField("Time", time.Since(startTime)).WithField("Bytes", len(buffer)).WithField("Region", updateRegion).Debug("Send frame")
	return
}

func confirmSend() (isOK bool, err error) {
	isOK = false
	buff := make([]byte, 10)
	len, err := serialPort.Read(buff)
	if err != nil {
		return
	} else if len == 0 {
		log.Warn("Timeout on confirmation")
		return
	} else if len > 1 {
		log.WithField("Buffer", buff[0:len]).Error("Multiple confirmation bytes recieved")
		return
	}

	switch buff[0] {
	case 'C':
		log.Error("Data corruption detected")
	case 'T':
		log.Error("Timeout on reciever")
	case 'O':
		log.Debug("Confirmation Recieved")
		isOK = true
	default:
		log.WithField("Code", buff[0]).Error("Unknown error code")
	}

	return
}

func getUpdateRegion(current image.Image, new image.Image) (shouldUpdate bool, region image.Rectangle) {
	if current == nil || current.Bounds() != new.Bounds() {
		return true, new.Bounds()
	}

	shouldUpdate, region.Min.X = findFirstDifferentPixel(current, new, true, true)
	if !shouldUpdate {
		return // Only the first one has to be checked
	}

	_, region.Min.Y = findFirstDifferentPixel(current, new, true, false)
	_, region.Max.X = findFirstDifferentPixel(current, new, false, true)
	_, region.Max.Y = findFirstDifferentPixel(current, new, false, false)

	region.Max = region.Max.Add(image.Pt(1, 1))

	return
}

func findFirstDifferentPixel(A image.Image, B image.Image, positive bool, isX bool) (differ bool, pos int) {
	outerMax := A.Bounds().Max.X
	innerMax := A.Bounds().Max.Y

	if !isX {
		outerMax, innerMax = innerMax, outerMax
	}

	for outer := 0; outer < outerMax; outer++ {
		for inner := 0; inner < innerMax; inner++ {
			pos := outer

			if !positive {
				pos = outerMax - outer - 1
			}

			x, y := pos, inner

			if !isX {
				x, y = y, x
			}

			ar, ag, ab, aa := A.At(x, y).RGBA()
			br, bg, bb, ba := B.At(x, y).RGBA()
			if ar != br || ag != bg || ab != bb || aa != ba {
				// log.
				// 	WithField("A", A.At(x, y)).
				// 	WithField("B", B.At(x, y)).
				// 	WithField("pos", image.Point{x, y}).
				// 	Info("Different Pixel")

				return true, pos
			}
		}
	}

	return false, 0
}

func encodeAndRotateImage(source image.Image, region image.Rectangle) []byte {
	if region != region.Canon() {
		log.WithField("Region", region).Fatal("Invalid region")
	}

	buff := make([]byte, region.Dx()*region.Dy()*2)

	for x := region.Min.X; x < region.Max.X; x++ {
		for y := region.Min.Y; y < region.Max.Y; y++ {
			bufferX := x - region.Min.X
			bufferY := y - region.Min.Y

			i := (bufferX*region.Dy() + bufferY) * 2 // TODO: Verify logic, look for cause of corruption
			sr, sg, sb, _ := source.At(x, y).RGBA()
			r := (byte)(sr >> 11)
			g := (byte)(sg >> 10)
			b := (byte)(sb >> 11)
			buff[i] = (byte)((g << 5) | b)
			buff[i+1] = (byte)((r << 3) | (g >> 3))
		}
	}

	if len(buff) != region.Size().X*region.Size().Y*2 {
		log.WithField("Length", len(buff)).WithField("Region", region).Fatal("invalid buffer size")
	}

	return buff
}

// func echo() {
// 	for {
// 		buff := make([]byte, 256)
// 		len, _ := serialPort.Read(buff)
// 		if len == 0 {
// 			continue
// 		}

// 		os.Stdout.Write(buff[0:len])
// 	}
// }

func InitSerial() (err error) {
	serialPort, err = serial.OpenPort(&serial.Config{
		Name:        "/dev/serial/by-id/usb-Arduino__www.arduino.cc__0043_64934333235351700140-if00",
		Baud:        1000000,
		ReadTimeout: 1000 * time.Millisecond,
	})

	if err != nil {
		return
	}

	err = serialPort.Flush()
	if err != nil {
		log.WithError(err).Fatal("Error flushing serial port")
	}

	// go echo()
	return
}
