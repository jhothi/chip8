package emulator

import (
	"math/rand"
	"time"
)

type cpu struct {
	v          [16]byte
	i          uint16
	pc         uint16
	delayTimer byte
	soundTimer byte
	stack      [16]uint16
	sp         uint16
	key        [16]byte
}

func (cpu *cpu) emulate(memory []uint8) {
	opCode := uint16(memory[cpu.pc]<<8 | memory[cpu.pc+1])
	switch opCode & 0xF000 {

	case 0x1000:
		//1nnn - JP addr
		//Jump to location nnn.
		//The interpreter sets the program counter to nnn.
		cpu.pc = opCode & 0x0FFF

	case 0x3000:
		//3xkk - SE Vx, byte
		//Skip next instruction if Vx = kk.
		//The interpreter compares register Vx to kk, and if they are equal, increments the program counter by 2.
		x := opCode & 0x0F00
		kk := byte(opCode & 0x00FF)
		if cpu.v[x] == kk {
			cpu.pc += 2
		}

	case 0x4000:
		//4xkk - SNE Vx, byte
		//Skip next instruction if Vx != kk.
		//The interpreter compares register Vx to kk, and if they are not equal, increments the program counter by 2.
		x := opCode & 0x0F00
		kk := byte(opCode & 0x00FF)
		if cpu.v[x] != kk {
			cpu.pc += 2
		}

	case 0x5000:
		//5xy0 - SE Vx, Vy
		//Skip next instruction if Vx = Vy.
		//The interpreter compares register Vx to register Vy, and if they are equal, increments the program counter by 2.
		x := opCode & 0x0F00
		y := opCode & 0x00F0
		if x == y {
			cpu.pc += 2
		}

	case 0x6000:
		//6xkk - LD Vx, byte
		//Set Vx = kk.
		//The interpreter puts the value kk into register Vx.
		x := opCode & 0x0F00
		kk := byte(opCode & 0x00FF)
		cpu.v[x] = kk

	case 0x7000:
		//7xkk - ADD Vx, byte
		//Set Vx = Vx + kk.
		//Adds the value kk to the value of register Vx, then stores the result in Vx.
		x := opCode & 0x0F00
		kk := byte(opCode & 0x00FF)
		cpu.v[x] += kk

	case 0x8000:
		x := opCode & 0x0F00
		y := opCode & 0x00F0
		switch opCode & 0x000F {

		case 0x0000:
			//8xy0 - LD Vx, Vy
			//Set Vx = Vy.
			//Stores the value of register Vy in register Vx.
			cpu.v[x] = cpu.v[y]

		case 0x0001:
			//8xy1 - OR Vx, Vy
			//Set Vx = Vx OR Vy.
			//Performs a bitwise OR on the values of Vx and Vy, then stores the result in Vx.
			cpu.v[x] = cpu.v[x] | cpu.v[y]

		case 0x0002:
			//8xy2 - AND Vx, Vy
			//Set Vx = Vx AND Vy.
			//Performs a bitwise AND on the values of Vx and Vy, then stores the result in Vx.
			cpu.v[x] = cpu.v[x] & cpu.v[y]

		case 0x0003:
			//8xy3 - XOR Vx, Vy
			//Set Vx = Vx XOR Vy.
			//Performs a bitwise exclusive OR on the values of Vx and Vy, then stores the result in Vx.
			cpu.v[x] = cpu.v[x] ^ cpu.v[y]

		case 0x0004:
			//8xy4 - ADD Vx, Vy
			//Set Vx = Vx + Vy, set VF = carry.
			//The values of Vx and Vy are added together. If the result is greater than 8 bits (i.e., > 255,)
			//VF is set to 1, otherwise 0. Only the lowest 8 bits of the result are kept, and stored in Vx.
			var result = uint16(cpu.v[x]) + uint16(cpu.v[y])
			if result > 255 {
				cpu.v[0xF] = 1
			} else {
				cpu.v[0xF] = 0
			}
			cpu.v[x] = byte(result & 0x00FF)

		case 0x0005:
			//8xy5 - SUB Vx, Vy
			//Set Vx = Vx - Vy, set VF = NOT borrow.
			//If Vx > Vy, then VF is set to 1, otherwise 0. Then Vy is subtracted from Vx, and the results stored in Vx.
			if cpu.v[x] > cpu.v[y] {
				cpu.v[0xF] = 1
			} else {
				cpu.v[0xF] = 0
			}
			cpu.v[x] -= cpu.v[y]

		case 0x0006:
			//8xy6 - SHR Vx {, Vy}
			//Set Vx = Vx SHR 1.
			//If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
			if cpu.v[x]&0x01 == 1 {
				cpu.v[0xF] = 1
			} else {
				cpu.v[0xF] = 0
			}
			cpu.v[x] = cpu.v[x] >> 1

		case 0x0007:
			//8xy7 - SUBN Vx, Vy
			//Set Vx = Vy - Vx, set VF = NOT borrow.
			//If Vy > Vx, then VF is set to 1, otherwise 0. Then Vx is subtracted from Vy, and the results stored in Vx.
			if cpu.v[y] > cpu.v[x] {
				cpu.v[0xF] = 1
			} else {
				cpu.v[0xF] = 0
			}
			cpu.v[x] = cpu.v[y] - cpu.v[x]

		default:
			//8xyE - SHL Vx {, Vy}
			//Set Vx = Vx SHL 1.
			//If the most-significant bit of Vx is 1, then VF is set to 1, otherwise to 0. Then Vx is multiplied by 2.
			if cpu.v[x]&0x80 == 1 {
				cpu.v[0xF] = 1
			} else {
				cpu.v[0xF] = 0
			}
			cpu.v[x] = cpu.v[x] << 1

		}

	case 0x9000:
		//9xy0 - SNE Vx, Vy
		//Skip next instruction if Vx != Vy.
		//The values of Vx and Vy are compared, and if they are not equal, the program counter is increased by 2.
		x := opCode & 0x0F00
		y := opCode & 0x00F0
		if x != y {
			cpu.pc += 2;
		}

	case 0xA000:
		//Annn - LD I, addr
		//Set I = nnn.
		//The value of register I is set to nnn.
		cpu.i = opCode & 0x0FFF
		cpu.pc += 2

	case 0xB000:
		//Bnnn - JP V0, addr
		//Jump to location nnn + V0.
		//The program counter is set to nnn plus the value of V0.
		cpu.pc = (opCode & 0x0FFF) + uint16(cpu.v[0])

	case 0xC000:
		//Cxkk - RND Vx, byte
		//Set Vx = random byte AND kk.
		//The interpreter generates a random number from 0 to 255, which is then ANDed with the value kk.
		//The results are stored in Vx.
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		x := opCode & 0x0F00
		kk := opCode & 0x00FF
		cpu.v[x] = byte(r1.Intn(255)) & byte(kk)

	case 0xD000:
		//Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision.
		//The interpreter reads n bytes from memory, starting at the address stored in I.
		//These bytes are then displayed as sprites on screen at coordinates (Vx, Vy). Sprites are XORed onto the existing screen.
		//If this causes any pixels to be erased, VF is set to 1, otherwise it is set to 0.
		//If the sprite is positioned so part of it is outside the coordinates of the display, it wraps around to the opposite side of the screen.
		//TODO: implement render

	case 0xE000:
		x := 0x0F00
		switch opCode & 0x00FF {

		case 0x9E:
			//Ex9E - SKP Vx
			//Skip next instruction if key with the value of Vx is pressed.
			//Checks the keyboard, and if the key corresponding to the value of Vx is currently in the down position, PC is increased by 2.
			if cpu.key[cpu.v[x]] == 1 { // TODO: fix check
				cpu.pc += 2
			}

		case 0xA1:
			//ExA1 - SKNP Vx
			//Skip next instruction if key with the value of Vx is not pressed.
			//Checks the keyboard, and if the key corresponding to the value of Vx is currently in the up position, PC is increased by 2.
			if cpu.key[cpu.v[x]] != 1 { // TODO: fix check
				cpu.pc += 2
			}

		}

	case 0xF000:
		x := opCode & 0x0F00
		switch opCode & 0x00FF {

		case 0x07:
			//Fx07 - LD Vx, DT
			//Set Vx = delay timer value.
			//The value of DT is placed into Vx.
			cpu.v[x] = cpu.delayTimer

		case 0x0A:
			//Fx0A - LD Vx, K
			//Wait for a key press, store the value of the key in Vx.
			//All execution stops until a key is pressed, then the value of that key is stored in Vx.
			//TODO read key press

		case 0x15:
			//Fx15 - LD DT, Vx
			//Set delay timer = Vx.
			//DT is set equal to the value of Vx.
			cpu.delayTimer = cpu.v[x]

		case 0x18:
			//Fx18 - LD ST, Vx
			//Set sound timer = Vx.
			//ST is set equal to the value of Vx.
			cpu.soundTimer = cpu.v[x]

		case 0x1E:
			//Fx1E - ADD I, Vx
			//Set I = I + Vx.
			//The values of I and Vx are added, and the results are stored in I.
			cpu.i += uint16(cpu.v[x])

		case 0x29:
			//Fx29 - LD F, Vx
			//Set I = location of sprite for digit Vx.
			//The value of I is set to the location for the hexadecimal sprite corresponding to the value of Vx.
			//TODO

		case 0x33:
			//Fx33 - LD B, Vx
			//Store BCD representation of Vx in memory locations I, I+1, and I+2.
			//The interpreter takes the decimal value of Vx, and places the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
			memory[cpu.i] = cpu.v[x] / 100
			memory[cpu.i+1] = (cpu.v[x] / 10) % 10
			memory[cpu.i+2] = (cpu.v[x] % 10) % 10

		case 0x55:
			//Fx55 - LD [I], Vx
			//Store registers V0 through Vx in memory starting at location I.
			//The interpreter copies the values of registers V0 through Vx into memory, starting at the address in I.
			for i, v := range cpu.v[0 : x+1] {
				memory[cpu.i+uint16(i)] = v
			}

		case 0x65:
			//Fx65 - LD Vx, [I]
			//Read registers V0 through Vx from memory starting at location I.
			//The interpreter reads values from memory starting at location I into registers V0 through Vx.
			for i := 0; uint16(i) <= x; i++ {
				cpu.v[i] = memory[cpu.i+uint16(i)]
			}

		}

	default:
		panic("Instruction not implemented")
	}
}

func newCpu() cpu {
	return cpu{pc: 0x200}
}
