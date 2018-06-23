package main

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"chip8/internal/emulator"
)

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 64, 32),
	}
	win, err := pixelgl.NewWindow(cfg)
	r := emulator.io{Window: win}
	r.Draw([]byte{0xF0,
		0x90,
		0xF0,
		0x90,
		0x90}, 0, 0)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Update()
	}
}

func main() {
	//filePath := os.Args[1]
	//emu := emulator.New()
	//data := emu.Start("/home/jagroop/go/src/chip8/roms/tetris.rom")
	//fmt.Println(data)
	pixelgl.Run(run)
}
