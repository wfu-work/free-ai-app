package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {
	const size = 64
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	drawLine(img, 32, 16, 32, 30, 5, white)
	drawLine(img, 32, 31, 22, 38, 5, white)
	drawLine(img, 32, 31, 42, 38, 5, white)
	drawLine(img, 22, 38, 32, 48, 5, white)
	drawLine(img, 42, 38, 32, 48, 5, white)

	drawCircle(img, 32, 16, 6, white)
	drawCircle(img, 32, 31, 5, white)
	drawCircle(img, 22, 38, 6, white)
	drawCircle(img, 42, 38, 6, white)
	drawCircle(img, 32, 48, 6, white)

	if err := os.MkdirAll("build/tray", 0o755); err != nil {
		panic(err)
	}
	file, err := os.Create("build/tray/freeai-template.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if err := png.Encode(file, img); err != nil {
		panic(err)
	}
}

func drawCircle(img *image.RGBA, cx, cy, radius int, col color.RGBA) {
	for y := cy - radius; y <= cy+radius; y++ {
		for x := cx - radius; x <= cx+radius; x++ {
			if (x-cx)*(x-cx)+(y-cy)*(y-cy) <= radius*radius {
				img.SetRGBA(x, y, col)
			}
		}
	}
}

func drawLine(img *image.RGBA, x1, y1, x2, y2, width int, col color.RGBA) {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	length := math.Hypot(dx, dy)
	if length == 0 {
		return
	}

	radius := float64(width) / 2
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			px := float64(x) + 0.5
			py := float64(y) + 0.5
			t := ((px-float64(x1))*dx + (py-float64(y1))*dy) / (length * length)
			if t < 0 {
				t = 0
			}
			if t > 1 {
				t = 1
			}
			nx := float64(x1) + t*dx
			ny := float64(y1) + t*dy
			if math.Hypot(px-nx, py-ny) <= radius {
				img.SetRGBA(x, y, col)
			}
		}
	}
}
