package main

import (
	"chip8/internal/emulator"
	"github.com/faiface/pixel/pixelgl"
)

func run() {
	emu := emulator.New()
	emu.Start("/home/jagroop/go/src/chip8/roms/tetris.rom")
}

func main() {
	//x := byte(8)
	//fmt.Println(x + byte(0xFF))
	pixelgl.Run(run)
}
