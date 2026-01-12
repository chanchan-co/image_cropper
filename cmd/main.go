package main

import (
	"flag"
	"log"

	"github.com/chanchan-co/image_cropper/internal/usecase"
)

func main() {
	inputDirPath := flag.String("input", "./tmp/images/input", "input directory path")
	outputDirPath := flag.String("output", "./tmp/images/output", "output directory path")
	cutPx := flag.Int("cut", 60, "pixels to cut from bottom")

	flag.Parse()

	if err := usecase.CropImages(*inputDirPath, *outputDirPath, *cutPx); err != nil {
		log.Fatalf("failed to crop images: %v", err)
	}
}
