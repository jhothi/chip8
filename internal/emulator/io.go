package emulator

import (
	"github.com/faiface/pixel/pixelgl"
	"image"
	"github.com/faiface/pixel"
)

type io struct {
	window *pixelgl.Window
	grid   [2048]byte
}

func (io *io) draw(sprite []byte, xPos byte, yPos byte) bool {
	collision := copySpriteToGrid(sprite, io.grid[:], xPos, yPos)
	toSprite(io.grid[:]).Draw(io.window, pixel.IM.Moved(io.window.Bounds().Center()))
	return collision
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
	//myfile, _ := os.Create("test.png")
	//png.Encode(myfile, img)
	pic := pixel.PictureDataFromImage(img)
	//fmt.Println(pic.Bounds())
	return pixel.NewSprite(pic, pic.Bounds())
}

func copySpriteToGrid(sprite []byte, grid []byte, xPos byte, yPos byte) bool {
	collision := false
	for index, value := range sprite {
		startPos := int(xPos) + ((int(yPos) + index) * 64)
		for bitIndex := 7; bitIndex >= 0; bitIndex-- {
			collision = collision || copyBitToGrid(getBit(value, byte(bitIndex)), grid, startPos+7-bitIndex)
		}
	}
	return collision
}

func copyBitToGrid(bit byte, grid []byte, pos int) bool {
	previousValue := grid[pos]
	grid[pos] ^= bit
	return previousValue == 1 && grid[pos] == 0
}

func getBit(num byte, bitIndex byte) byte {
	return num & (1 << bitIndex) >> byte(bitIndex)
}