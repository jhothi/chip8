package main

import (
	"github.com/faiface/pixel/pixelgl"
	"chip8/internal/emulator"
)

func run() {
	emu := emulator.New()
	emu.Start("/home/jagroop/go/src/chip8/roms/tetris.rom")
}

func main() {
	pixelgl.Run(run)
}
