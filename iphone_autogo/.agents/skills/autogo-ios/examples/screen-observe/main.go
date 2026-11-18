package main

import (
	"fmt"

	"github.com/Dasongzi1366/AutoGo/files"
	"github.com/Dasongzi1366/AutoGo/images"
	"github.com/Dasongzi1366/AutoGo/ppocr"
)

func main() {
	img := images.CaptureScreen(0, 0, 0, 0)
	if img == nil {
		fmt.Println("screen capture failed")
		return
	}

	path := files.Path("screen-observe.png")
	if path != "" && !images.Save(img, path, 100) {
		fmt.Println("screenshot save failed")
	}
	if path != "" {
		fmt.Println("screenshot:", path)
	}

	ocr := ppocr.New("v5")
	if ocr == nil {
		fmt.Println("PPOCR initialization failed")
		return
	}
	defer ocr.Close()

	for _, result := range ocr.Ocr(0, 0, 0, 0, "") {
		fmt.Printf("[%d,%d] score=%.2f %s\n", result.CenterX, result.CenterY, result.Score, result.Label)
	}
}
