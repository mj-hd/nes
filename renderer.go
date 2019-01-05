package main

import "image/color"

type Renderer interface {
	SetPixel(x, y int, col color.RGBA)
	Render()
}
