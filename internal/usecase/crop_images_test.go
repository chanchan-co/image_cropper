package usecase

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildImagePath(t *testing.T) {
	tests := []struct {
		name     string
		dirPath  string
		fileName string
		want     string
	}{
		{
			name:     "normal path",
			dirPath:  "/tmp/images",
			fileName: "test.jpg",
			want:     "/tmp/images/test.jpg",
		},
		{
			name:     "with trailing slash",
			dirPath:  "/tmp/images/",
			fileName: "test.png",
			want:     "/tmp/images/test.png",
		},
		{
			name:     "relative path",
			dirPath:  "./images",
			fileName: "test.jpg",
			want:     "images/test.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildImagePath(tt.dirPath, tt.fileName)
			if got != tt.want {
				t.Errorf("buildImagePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeImg(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("decode jpeg", func(t *testing.T) {
		imgPath := filepath.Join(tmpDir, "test.jpg")
		img := createTestImage(100, 100)
		file, err := os.Create(imgPath)
		if err != nil {
			t.Fatal(err)
		}
		jpeg.Encode(file, img, nil)
		file.Close()

		decoded, format, err := decodeImg(imgPath)
		if err != nil {
			t.Fatalf("decodeImg() error = %v", err)
		}
		if format != "jpeg" {
			t.Errorf("format = %v, want jpeg", format)
		}
		if decoded == nil {
			t.Error("decoded image is nil")
		}
	})

	t.Run("decode png", func(t *testing.T) {
		imgPath := filepath.Join(tmpDir, "test.png")
		img := createTestImage(100, 100)
		file, err := os.Create(imgPath)
		if err != nil {
			t.Fatal(err)
		}
		png.Encode(file, img)
		file.Close()

		decoded, format, err := decodeImg(imgPath)
		if err != nil {
			t.Fatalf("decodeImg() error = %v", err)
		}
		if format != "png" {
			t.Errorf("format = %v, want png", format)
		}
		if decoded == nil {
			t.Error("decoded image is nil")
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, _, err := decodeImg(filepath.Join(tmpDir, "notfound.jpg"))
		if err == nil {
			t.Error("expected error for non-existent file")
		}
	})

	t.Run("invalid image file", func(t *testing.T) {
		imgPath := filepath.Join(tmpDir, "invalid.jpg")
		os.WriteFile(imgPath, []byte("not an image"), 0644)

		_, _, err := decodeImg(imgPath)
		if err == nil {
			t.Error("expected error for invalid image file")
		}
	})
}

func TestSaveImage(t *testing.T) {
	tmpDir := t.TempDir()
	img := createTestImage(50, 50)

	t.Run("save jpeg", func(t *testing.T) {
		imgPath := filepath.Join(tmpDir, "output.jpg")
		err := saveImage(imgPath, img, "jpeg")
		if err != nil {
			t.Fatalf("saveImage() error = %v", err)
		}

		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			t.Error("saved file does not exist")
		}
	})

	t.Run("save png", func(t *testing.T) {
		imgPath := filepath.Join(tmpDir, "output.png")
		err := saveImage(imgPath, img, "png")
		if err != nil {
			t.Fatalf("saveImage() error = %v", err)
		}

		if _, err := os.Stat(imgPath); os.IsNotExist(err) {
			t.Error("saved file does not exist")
		}
	})

	t.Run("unsupported format", func(t *testing.T) {
		imgPath := filepath.Join(tmpDir, "output.bmp")
		err := saveImage(imgPath, img, "bmp")
		if err == nil {
			t.Error("expected error for unsupported format")
		}
	})

	t.Run("invalid directory", func(t *testing.T) {
		imgPath := filepath.Join(tmpDir, "nonexistent", "output.jpg")
		err := saveImage(imgPath, img, "jpeg")
		if err == nil {
			t.Error("expected error for invalid directory")
		}
	})
}

func TestCropImages(t *testing.T) {
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	os.MkdirAll(inputDir, 0755)

	t.Run("successful crop", func(t *testing.T) {
		img := createTestImage(100, 200)
		jpegPath := filepath.Join(inputDir, "test.jpg")
		pngPath := filepath.Join(inputDir, "test.png")

		saveTestImage(t, jpegPath, img, "jpeg")
		saveTestImage(t, pngPath, img, "png")

		err := CropImages(inputDir, outputDir, 50)
		if err != nil {
			t.Fatalf("CropImages() error = %v", err)
		}

		croppedJpeg := filepath.Join(outputDir, "cropped_test.jpg")
		croppedPng := filepath.Join(outputDir, "cropped_test.png")

		if _, err := os.Stat(croppedJpeg); os.IsNotExist(err) {
			t.Error("cropped jpeg file does not exist")
		}
		if _, err := os.Stat(croppedPng); os.IsNotExist(err) {
			t.Error("cropped png file does not exist")
		}
	})

	t.Run("input directory not found", func(t *testing.T) {
		err := CropImages(filepath.Join(tmpDir, "notfound"), outputDir, 50)
		if err == nil {
			t.Error("expected error for non-existent input directory")
		}
	})

	t.Run("creates output directory", func(t *testing.T) {
		newOutputDir := filepath.Join(tmpDir, "new_output")
		img := createTestImage(100, 100)
		testImgPath := filepath.Join(inputDir, "create_test.jpg")
		saveTestImage(t, testImgPath, img, "jpeg")

		err := CropImages(inputDir, newOutputDir, 10)
		if err != nil {
			t.Fatalf("CropImages() error = %v", err)
		}

		if _, err := os.Stat(newOutputDir); os.IsNotExist(err) {
			t.Error("output directory was not created")
		}
	})

	t.Run("skip subdirectories", func(t *testing.T) {
		subDir := filepath.Join(inputDir, "subdir")
		os.MkdirAll(subDir, 0755)

		outputDir2 := filepath.Join(tmpDir, "output2")
		err := CropImages(inputDir, outputDir2, 50)
		if err != nil {
			t.Fatalf("CropImages() error = %v", err)
		}

		entries, _ := os.ReadDir(outputDir2)
		for _, entry := range entries {
			if entry.IsDir() {
				t.Error("subdirectories should be skipped")
			}
		}
	})
}

func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	return img
}

func saveTestImage(t *testing.T, path string, img image.Image, format string) {
	t.Helper()
	file, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	switch format {
	case "jpeg":
		jpeg.Encode(file, img, nil)
	case "png":
		png.Encode(file, img)
	}
}
