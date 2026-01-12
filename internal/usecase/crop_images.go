package usecase

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"

	"github.com/chanchan-co/image_cropper/internal/transform"
)

func CropImages(inputDirPath, outputDirPath string, cutPx int) error {
	if err := os.MkdirAll(outputDirPath, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	entries, err := os.ReadDir(inputDirPath)
	if err != nil {
		return fmt.Errorf("failed to read input directory: %w", err)
	}

	successCount := 0
	failCount := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		inputImgPath := buildImagePath(inputDirPath, entry.Name())
		outputImgPath := buildImagePath(outputDirPath, "cropped_"+entry.Name())

		img, format, err := decodeImg(inputImgPath)
		if err != nil {
			log.Printf("failed to decode %s: %v", entry.Name(), err)
			failCount++
			continue
		}

		croppedImg := transform.CropBottom(img, cutPx)
		if err := saveImage(outputImgPath, croppedImg, format); err != nil {
			log.Printf("failed to save %s: %v", entry.Name(), err)
			failCount++
			continue
		}

		log.Printf("processed: %s", entry.Name())
		successCount++
	}

	log.Printf("completed: %d succeeded, %d failed", successCount, failCount)
	return nil
}

func buildImagePath(dirPath, fileName string) string {
	return filepath.Join(dirPath, fileName)
}

func decodeImg(imgPath string) (image.Image, string, error) {
	imgFile, err := os.Open(imgPath)
	if err != nil {
		return nil, "", err
	}
	defer imgFile.Close()

	img, format, err := image.Decode(imgFile)
	if err != nil {
		return nil, "", err
	}

	if format != "jpeg" && format != "png" {
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}

	return img, format, nil
}

func saveImage(imgPath string, img image.Image, format string) error {
	imgFile, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	defer imgFile.Close()

	switch format {
	case "jpeg":
		return jpeg.Encode(imgFile, img, nil)
	case "png":
		return png.Encode(imgFile, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
