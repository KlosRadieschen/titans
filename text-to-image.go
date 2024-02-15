package main

import (
	"github.com/fogleman/gg"
)

func drawText(text []string) {
	const W = 800
	const H = 800
	dc := gg.NewContext(W, H)
	dc.SetRGB(50, 50, 50)
	dc.Clear()
	dc.SetRGB(255, 0, 0)
	dc.LoadFontFace("/home/Nicolas/go-workspace/src/titans/OpenSans-Bold.ttf", 30)
	const h = 24
	for i, line := range text {
		y := H/2 - h*len(text)/2 + i*h
		dc.DrawStringAnchored(line, 400, float64(y), 0.5, 0.5)
	}
	dc.SavePNG(directory + "table.png")
}
