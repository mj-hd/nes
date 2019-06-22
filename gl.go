package main

import (
	"image"
	"image/color"
	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	runtime.LockOSThread()
}

type GLRenderer struct {
	window  *glfw.Window
	image   *image.RGBA
	texture uint32
}

func NewGLRenderer(window *glfw.Window) Renderer {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return &GLRenderer{
		window: window,
		image: &image.RGBA{

			Pix: make([]uint8, 4*256*240), // TODO: remove magic number
			Rect: image.Rectangle{
				Min: image.Point{
					0, 0,
				},
				Max: image.Point{
					256, 240,
				},
			},
		},
		texture: texture,
	}
}

func (g *GLRenderer) SetPixel(x, y int, col color.RGBA) {
	base := y*256*4 + x*4
	g.image.Pix[base] = col.R
	g.image.Pix[base+1] = col.G
	g.image.Pix[base+2] = col.B
	g.image.Pix[base+3] = col.A
}

func (g *GLRenderer) Render() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	s := g.image.Rect.Size()

	gl.BindTexture(gl.TEXTURE_2D, g.texture)

	gl.TexImage2D(
		gl.TEXTURE_2D, 0,
		gl.RGBA,
		int32(s.X), int32(s.Y), 0,
		gl.RGBA,
		gl.UNSIGNED_BYTE, gl.Ptr(g.image.Pix),
	)

	w, h := g.window.GetFramebufferSize()
	s1 := float32(w) / 256
	s2 := float32(h) / 240
	f := float32(1)
	var x, y float32
	if s1 >= s2 {
		x = f * s2 / s1
		y = f
	} else {
		x = f
		y = f * s1 / s2
	}
	gl.Begin(gl.QUADS)
	gl.TexCoord2f(0, 1)
	gl.Vertex2f(-x, -y)
	gl.TexCoord2f(1, 1)
	gl.Vertex2f(x, -y)
	gl.TexCoord2f(1, 0)
	gl.Vertex2f(x, y)
	gl.TexCoord2f(0, 0)
	gl.Vertex2f(-x, y)
	gl.End()

	gl.BindTexture(gl.TEXTURE_2D, 0)
	g.window.SwapBuffers()
}
