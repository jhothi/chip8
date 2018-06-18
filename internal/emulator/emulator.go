package emulator

import "io/ioutil"

type emulator struct {
	cpu    cpu
	memory [4096]byte
}

func New() emulator  {
	cpu  := newCpu()
	return emulator{cpu:cpu}
}

func (e *emulator) Start(romPath string) [4096]byte {
	e.loadRom(romPath)
	e.cpu.emulate(e.memory[:])
	return e.memory
}

func (e *emulator) loadRom(romPath string) [4096]byte {
	data, _ := ioutil.ReadFile(romPath)
	for i, v := range data {
		e.memory[i+512] = v
	}
	return e.memory
}
