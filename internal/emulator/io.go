package emulator

import (
	"github.com/faiface/pixel/pixelgl"
	"image"
	"github.com/faiface/pixel"
	"fmt"
	"github.com/go-gl/glfw/v3.2/glfw"
	"errors"
	"os"
	"image/png"
)

var keyMap = map[pixelgl.Button]byte{
	pixelgl.Button(glfw.Key1): 1,
	pixelgl.Button(glfw.Key2): 2,
	pixelgl.Button(glfw.Key3): 3,
	pixelgl.Button(glfw.Key4): 0xC,
	pixelgl.Button(glfw.KeyQ): 4,
	pixelgl.Button(glfw.KeyW): 5,
	pixelgl.Button(glfw.KeyE): 6,
	pixelgl.Button(glfw.KeyR): 0xD,
	pixelgl.Button(glfw.KeyA): 7,
	pixelgl.Button(glfw.KeyS): 8,
	pixelgl.Button(glfw.KeyD): 9,
	pixelgl.Button(glfw.KeyF): 0xE,
	pixelgl.Button(glfw.KeyZ): 0xA,
	pixelgl.Button(glfw.KeyX): 0,
	pixelgl.Button(glfw.KeyC): 0xB,
	pixelgl.Button(glfw.KeyV): 0xF,
}

type io struct {
	window *pixelgl.Window
	grid   [2048]byte
	keys   [16]byte
}

func (io *io) draw(sprite []byte, xPos byte, yPos byte) bool {
	collision := copySpriteToGrid(sprite, io.grid[:], xPos, yPos)
	fmt.Println(io.grid)
	toSprite(io.grid[:]).Draw(io.window, pixel.IM.Scaled(pixel.Vec{0,0}, 16).Moved(io.window.Bounds().Center()))
	return collision
}


func (io *io) readKeyPress() (byte, error) {
	for k, v := range keyMap {
		if io.window.JustPressed(k) {
			io.keys[v] = 1
			return v, nil
		}
	}
	return 0, errors.New("no key pressed")
}

func (io *io) getKeyPress(pos byte) byte {
	if io.keys[pos] == 1 {
		io.keys[pos] = 0
		return 1
	}
	return io.keys[pos]
}

func (io *io) clearDisplay() {
	for i := 0; i < len(io.grid); i++ {
		io.grid[i] = 0
	}
}

func toSprite(grid []byte) *pixel.Sprite {
	img := image.NewGray(image.Rect(0, 0, 64, 32))
	for i, v := range grid {
		if v == 1 {
			img.Pix[i] = 255
		} else {
			img.Pix[i] = 0
		}
	}
	//fmt.Println(grid)
	//fmt.Println(img.Pix)
	myfile, _ := os.Create("test.png")
	png.Encode(myfile, img)
	pic := pixel.PictureDataFromImage(img)
	//fmt.Println(pic.Bounds())
	return pixel.NewSprite(pic, pic.Bounds())
}

func copySpriteToGrid(sprite []byte, grid []byte, xPos byte, yPos byte) bool {
	collision := false
	for index, value := range sprite {
		startPos := int(xPos) + ((int(yPos) + index) * 64) % 2048
		for bitIndex := 7; bitIndex >= 0; bitIndex-- {
			collision = copyBitToGrid(getBit(value, byte(bitIndex)), grid, startPos+7-bitIndex) || collision
		}
	}
	return collision
}

func copyBitToGrid(bit byte, grid []byte, pos int) bool {
	fmt.Printf("pos %v", pos)
	previousValue := grid[pos]
	grid[pos] ^= bit
	return previousValue == 1 && grid[pos] == 0
}

func getBit(num byte, bitIndex byte) byte {
	return num & (1 << bitIndex) >> byte(bitIndex)
}

func newIo() io {
	cfg := pixelgl.WindowConfig{
		Title:  "Chip 8",
		Bounds: pixel.R(0, 0, 1024, 512),
		VSync:  true,
	}
	win, _ := pixelgl.NewWindow(cfg)
	return io{window: win}
}
