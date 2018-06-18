package main

import (
	"os"
	"chip8/internal/emulator"
	"fmt"
)

func main() {
	filePath := os.Args[1]
	emu := emulator.New()
	data := emu.Start(filePath)
	fmt.Println(data)
}
