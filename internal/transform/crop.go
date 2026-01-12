package transform

import (
	"image"
	"log"
)

func CropBottom(img image.Image, cutPx int) (image.Image) {
	b := img.Bounds()

	if cutPx <= 0 {
		return img
	}

	if b.Dy() <= cutPx {
		return img
	}

	rect := image.Rect(
		b.Min.X,
		b.Min.Y,
		b.Max.X,
		b.Max.Y - cutPx,
	)

	sub, ok := img.(interface {
		SubImage(r image.Rectangle) image.Image
	})

	if !ok {
		log.Printf("warning: image type does not support SubImage, returning original image")
		return img
	}

	return sub.SubImage(rect)
}
