package transform

import (
	"image"
	"image/color"
	"testing"
)

func TestCropBottom(t *testing.T) {
	tests := []struct {
		name     string
		img      image.Image
		cutPx    int
		wantH    int
		wantW    int
	}{
		{
			name:  "normal crop",
			img:   image.NewRGBA(image.Rect(0, 0, 100, 200)),
			cutPx: 50,
			wantH: 150,
			wantW: 100,
		},
		{
			name:  "cutPx is zero",
			img:   image.NewRGBA(image.Rect(0, 0, 100, 200)),
			cutPx: 0,
			wantH: 200,
			wantW: 100,
		},
		{
			name:  "cutPx is negative",
			img:   image.NewRGBA(image.Rect(0, 0, 100, 200)),
			cutPx: -10,
			wantH: 200,
			wantW: 100,
		},
		{
			name:  "cutPx equals image height",
			img:   image.NewRGBA(image.Rect(0, 0, 100, 200)),
			cutPx: 200,
			wantH: 200,
			wantW: 100,
		},
		{
			name:  "cutPx greater than image height",
			img:   image.NewRGBA(image.Rect(0, 0, 100, 200)),
			cutPx: 300,
			wantH: 200,
			wantW: 100,
		},
		{
			name:  "small crop",
			img:   image.NewRGBA(image.Rect(0, 0, 50, 50)),
			cutPx: 10,
			wantH: 40,
			wantW: 50,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CropBottom(tt.img, tt.cutPx)
			bounds := result.Bounds()
			gotH := bounds.Dy()
			gotW := bounds.Dx()

			if gotH != tt.wantH {
				t.Errorf("CropBottom() height = %v, want %v", gotH, tt.wantH)
			}
			if gotW != tt.wantW {
				t.Errorf("CropBottom() width = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}

func TestCropBottom_UnsupportedImageType(t *testing.T) {
	unsupportedImg := &unsupportedImage{
		bounds: image.Rect(0, 0, 100, 200),
	}

	result := CropBottom(unsupportedImg, 50)

	if result != unsupportedImg {
		t.Error("CropBottom() should return original image for unsupported type")
	}
}

type unsupportedImage struct {
	bounds image.Rectangle
}

func (u *unsupportedImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (u *unsupportedImage) Bounds() image.Rectangle {
	return u.bounds
}

func (u *unsupportedImage) At(x, y int) color.Color {
	return color.RGBA{R: 0, G: 0, B: 0, A: 255}
}
