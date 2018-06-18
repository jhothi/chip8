package main

import (
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
)


func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
	}
	win, err := pixelgl.NewWindow(cfg)
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
	//data := emu.Start(filePath)
	//fmt.Println(data)
	pixelgl.Run(run)
}
