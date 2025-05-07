package utils

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
)

func CreateSparkline(data []float64, width, height int, backColor color.Color, foreColor color.Color) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{backColor}, image.Point{}, draw.Src)

	if len(data) == 0 {
		return nil, nil
	}

	min, max := data[0], data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	scale := func(value float64) float64 {
		if max == min {
			return float64(height) / 2
		}
		return (value - min) / (max - min) * float64(height)
	}

	for i := 0; i < len(data)-1; i++ {
		x1 := float64(i) / float64(len(data)-1) * float64(width)
		y1 := float64(height) - scale(data[i])
		x2 := float64(i+1) / float64(len(data)-1) * float64(width)
		y2 := float64(height) - scale(data[i+1])

		drawLine(img, int(x1), int(y1), int(x2), int(y2), foreColor)
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, col color.Color) {
	dx := math.Abs(float64(x2 - x1))
	dy := math.Abs(float64(y2 - y1))
	sx := -1
	sy := -1

	if x1 < x2 {
		sx = 1
	}
	if y1 < y2 {
		sy = 1
	}

	err := dx - dy
	for {
		img.Set(x1, y1, col)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}
