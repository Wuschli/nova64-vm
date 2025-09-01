package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"maps"
	"slices"

	"github.com/rotisserie/eris"
)

var instructionMap = map[byte]func(*nova64Cpu, uint32) error{
	0x00: noop,  // NOOP
	0x01: push,  // PUSH
	0x02: drop,  // DROP
	0x03: dup,   // DUP
	0x04: swap,  // SWAP
	0x10: add,   // ADD
	0x11: sub,   // SUB
	0x12: mul,   // MUL
	0x13: div,   // DIV
	0x14: mod,   // MOD
	0x20: and,   // AND
	0x21: or,    // OR
	0x22: xor,   // XOR
	0x23: not,   // NOT
	0x30: cmp,   // CMP
	0x40: jmp,   // JMP
	0x41: call,  // CALL
	0x42: ret,   // RET
	0x43: jmpz,  // JMPZ
	0x44: jmplt, // JMPLT
	0x45: jmpgt, // JMPGT
	0x50: fetch, // FETCH
	0x51: store, // STORE
	0x60: in,    // IN
	0x61: out,   // OUT
	0x70: spawn, // SPAWN
	0x71: yield, // YIELD
	0x72: wait,  // WAIT
	0xFF: kill,  // KILL
}

type nova64Cpu struct {
	Tasks                      []*nova64Task
	Ports                      map[uint32]Port
	Ram                        []uint32
	Labels                     map[string]uint32
	StackSize                  uint32
	TicksPerTask               int32
	Halted                     bool
	activeTaskIndex            int
	ticksRemainingOnActiveTask int32
	usedStackOffsets           map[*nova64Task]uint32
	nextTaskId                 uint32
}

type nova64Task struct {
	ID, IP, SP, Flags, Waiting, stackBase uint32
	X                                     int32
	killed                                bool
}

type Port interface {
	Read() uint32
	Write(value uint32)
	HasData() bool
}

func NewCpu(memory uint32) nova64Cpu {
	cpu := nova64Cpu{
		Tasks:            make([]*nova64Task, 0),
		Ports:            make(map[uint32]Port),
		Ram:              make([]uint32, memory/4),
		StackSize:        512,
		TicksPerTask:     512,
		usedStackOffsets: make(map[*nova64Task]uint32, 1),
	}
	return cpu
}

func (cpu *nova64Cpu) Load(data []byte) error {
	labels, err := parseImage(data, cpu.Ram[:])
	if err != nil {
		return err
	}
	cpu.Labels = labels
	cpu.SpawnTask(0)
	return nil
}

func (cpu *nova64Cpu) Tick() {
	if cpu.Halted || cpu.doTaskScheduling() {
		return
	}
	task := cpu.ActiveTask()

	if cpu.ActiveTask().Waiting != 0 {
		return
	}

	op := cpu.Ram[task.IP]
	opCode := byte(op & 0xFF)
	operand := (op >> 8) & 0x00FFFFFF
	instructionFunc := instructionMap[opCode]
	err := instructionFunc(cpu, operand)
	if err != nil {
		log.Fatalln(eris.ToString(err, true))
		cpu.ActiveTask().killed = true
		// panic(err)
	}

	task.IP++
}

func (cpu *nova64Cpu) doTaskScheduling() bool {
	cpu.ticksRemainingOnActiveTask--
	if cpu.ticksRemainingOnActiveTask < 0 || (cpu.ActiveTask() != nil && (cpu.ActiveTask().Waiting != 0 || cpu.ActiveTask().killed)) {

		if cpu.ActiveTask() != nil && cpu.ActiveTask().killed {
			delete(cpu.usedStackOffsets, cpu.ActiveTask())
			cpu.Tasks = slices.Delete(cpu.Tasks, cpu.activeTaskIndex, cpu.activeTaskIndex+1)
		} else {
			cpu.activeTaskIndex++
			cpu.activeTaskIndex %= len(cpu.Tasks)
			cpu.ticksRemainingOnActiveTask = cpu.TicksPerTask
		}

		if len(cpu.Tasks) == 0 {
			cpu.Halted = true
			return true
		}
	}
	return false
}

func (cpu *nova64Cpu) ActiveTask() *nova64Task {
	return cpu.Tasks[cpu.activeTaskIndex]
}

func (cpu *nova64Cpu) SpawnTask(ip uint32) {
	usedOffsets := slices.Sorted(maps.Values(cpu.usedStackOffsets))
	offset := uint32(0)

	for _, i := range usedOffsets {
		if offset == i {
			offset = i + 1
		}
	}

	task := nova64Task{
		ID:        cpu.nextTaskId,
		IP:        ip,
		stackBase: uint32(len(cpu.Ram)) - uint32(cpu.StackSize)*(offset+1),
		SP:        uint32(len(cpu.Ram)) - uint32(cpu.StackSize)*(offset+1) - 1,
	}

	cpu.nextTaskId++
	cpu.Tasks = append(cpu.Tasks, &task)
	cpu.usedStackOffsets[&task] = offset
}

func (cpu *nova64Cpu) push(value uint32) error {
	cpu.ActiveTask().SP++
	if cpu.ActiveTask().SP >= cpu.ActiveTask().stackBase+cpu.StackSize {
		return errors.New("StackOverflow")
	}
	cpu.Ram[cpu.ActiveTask().SP] = value
	return nil
}

func (cpu *nova64Cpu) drop() error {
	if cpu.ActiveTask().SP < cpu.ActiveTask().stackBase {
		return errors.New("StackUnderflow")
	}
	cpu.ActiveTask().SP--
	return nil
}

func (cpu *nova64Cpu) pop() (uint32, error) {
	if cpu.ActiveTask().SP < cpu.ActiveTask().stackBase {
		return 0, errors.New("StackUnderflow")
	}
	value := cpu.Ram[cpu.ActiveTask().SP]
	cpu.ActiveTask().SP--
	return value, nil
}

func (cpu *nova64Cpu) dup(offset uint32) error {
	if cpu.ActiveTask().SP-offset < cpu.ActiveTask().stackBase {
		return errors.New("StackUnderflow")
	}
	return cpu.push(cpu.Ram[cpu.ActiveTask().stackBase-1-offset])
}

func (cpu *nova64Cpu) swap() error {
	if cpu.ActiveTask().SP-1 < cpu.ActiveTask().stackBase {
		return errors.New("StackUnderflow")
	}
	temp := cpu.Ram[cpu.ActiveTask().SP]
	cpu.Ram[cpu.ActiveTask().SP] = cpu.Ram[cpu.ActiveTask().SP-1]
	cpu.Ram[cpu.ActiveTask().SP-1] = temp
	return nil
}

func (cpu *nova64Cpu) yield() {
	cpu.ticksRemainingOnActiveTask = 0
}

func parseImage(data []byte, ram []uint32) (map[string]uint32, error) {
	reader := bytes.NewReader(data)
	var labelCount uint32
	if err := binary.Read(reader, binary.LittleEndian, &labelCount); err != nil {
		return nil, err
	}

	labels := make(map[string]uint32)

	var label bytes.Buffer
	var address uint32
	for i := 0; i < int(labelCount); i++ {
		// read address
		if err := binary.Read(reader, binary.LittleEndian, &address); err != nil {
			return nil, err
		}

		// read string until NULL byte
		for {
			b, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			if b == 0 {
				break
			}
			if err = label.WriteByte(b); err != nil {
				return nil, err
			}
		}
		labels[label.String()] = address
	}

	// copy rest of the buffer to ram
	imageSize := reader.Len() / 4
	for i := range imageSize {
		if err := binary.Read(reader, binary.LittleEndian, &ram[i]); err != nil {
			return nil, err
		}
	}
	return labels, nil
}
