package main

import (
	"os"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {

	if err := glfw.Init(); err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(300, 300, "nes", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		panic(err)
	}
	gl.Enable(gl.TEXTURE_2D)

	r := NewGLRenderer(window)

	nes, err := NewNES(os.Args[1], r)
	if err != nil {
		panic(err)
	}

	//fmt.Println("======CPU=======")
	//for i := uint16(0); i < 0xFFFF; i++ {
	//	if i%16 == 0 {
	//		fmt.Printf("\n%04x: ", i)
	//	}
	//	fmt.Printf("%04x ", nes.CPU.get(i))
	//}
	//fmt.Println("======PPU=======")
	//for i := uint16(0); i < 0x4000; i++ {
	//	if i%16 == 0 {
	//		fmt.Printf("\n%04x: ", i)
	//	}
	//	fmt.Printf("%04x ", nes.PPU.get(i))
	//}

	nes.PowerOn()

	for !window.ShouldClose() {
		nes.Tick()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
