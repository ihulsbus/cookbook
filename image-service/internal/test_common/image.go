package testcommon

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
)

func createImage() *image.RGBA {
	width := 200
	height := 100

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	cyan := color.RGBA{100, 200, 200, 0xff}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2:
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2:
				img.Set(x, y, color.White)
			default:
			}
		}
	}

	return img
}

func CreateFile() (bytes.Buffer, error) {
	var buff bytes.Buffer
	var err error

	img := createImage()

	err = png.Encode(&buff, img)
	if err != nil {
		return buff, err
	}

	return buff, nil

}
