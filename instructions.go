package main

import (
	"errors"
	"fmt"

	"github.com/rotisserie/eris"

	comp "cmp"
)

const TOP uint32 = 0xFFFFFF

func noop(cpu *nova64Cpu, operand uint32) error {
	return nil
}

func push(cpu *nova64Cpu, operand uint32) error {
	if err := cpu.push(operand); err != nil {
		return eris.Wrapf(err, "instruction PUSH failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func drop(cpu *nova64Cpu, operand uint32) error {
	if err := cpu.drop(); err != nil {
		return eris.Wrapf(err, "instruction DROP failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func dup(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return eris.Wrapf(err, "instruction DUP failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
		}
	}
	if err := cpu.dup(operand); err != nil {
		return eris.Wrapf(err, "instruction DUP failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func swap(cpu *nova64Cpu, operand uint32) error {
	if err := cpu.swap(); err != nil {
		return eris.Wrapf(err, "instruction SWAP failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func add(cpu *nova64Cpu, operand uint32) error {
	b, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction ADD failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	a, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction ADD failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	err = cpu.push(uint32(int32(a) + int32(b)))
	if err != nil {
		return eris.Wrapf(err, "instruction ADD failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func sub(cpu *nova64Cpu, operand uint32) error {
	b, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction SUB failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	a, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction SUB failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	err = cpu.push(uint32(int32(a) - int32(b)))
	if err != nil {
		return eris.Wrapf(err, "instruction SUB failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func mul(cpu *nova64Cpu, operand uint32) error {
	b, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction MUL failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	a, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction MUL failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	err = cpu.push(uint32(int32(a) * int32(b)))
	if err != nil {
		return eris.Wrapf(err, "instruction MUL failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func div(cpu *nova64Cpu, operand uint32) error {
	b, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction DIV failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	a, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction DIV failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	err = cpu.push(uint32(int32(a) / int32(b)))
	if err != nil {
		return eris.Wrapf(err, "instruction DIV failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func mod(cpu *nova64Cpu, operand uint32) error {
	b, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction MOD failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	a, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction MOD failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	err = cpu.push(uint32(int32(a) % int32(b)))
	if err != nil {
		return eris.Wrapf(err, "instruction MOD failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func and(cpu *nova64Cpu, operand uint32) error {
	b, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction AND failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	a, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction AND failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	err = cpu.push(a & b)
	if err != nil {
		return eris.Wrapf(err, "instruction AND failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func or(cpu *nova64Cpu, operand uint32) error {
	b, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction OR failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	a, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction OR failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	err = cpu.push(a | b)
	if err != nil {
		return eris.Wrapf(err, "instruction OR failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func xor(cpu *nova64Cpu, operand uint32) error {
	b, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction XOR failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	a, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction XOR failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	err = cpu.push(a ^ b)
	if err != nil {
		return eris.Wrapf(err, "instruction XOR failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func not(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction XOR failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	err = cpu.push(^a)
	if err != nil {
		return eris.Wrapf(err, "instruction XOR failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func cmp(cpu *nova64Cpu, operand uint32) error {
	b, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction CMP failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	a, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction CMP failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}

	cpu.ActiveTask().X = int32(comp.Compare(int32(a), int32(b)))
	return nil
}

func jmp(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return eris.Wrapf(err, "instruction JMP failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
		}
	}
	cpu.ActiveTask().IP = operand - 1
	// fmt.Printf("Jump to %#08x\n", operand)
	return nil
}

func call(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return eris.Wrapf(err, "instruction CALL failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
		}
	}

	err := cpu.push(cpu.ActiveTask().IP)
	if err != nil {
		return eris.Wrapf(err, "instruction CALL failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}

	cpu.ActiveTask().IP = operand - 1
	return nil
}

func ret(cpu *nova64Cpu, operand uint32) error {
	value, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction RET failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	cpu.ActiveTask().IP = value - 1
	return nil
}

func jmpz(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return eris.Wrapf(err, "instruction JMPZ failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
		}
	}
	if cpu.ActiveTask().X == 0 {
		cpu.ActiveTask().IP = operand - 1
		// fmt.Printf("Jump to %#08x\n", operand)
	}
	return nil
}

func jmplt(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return eris.Wrapf(err, "instruction JMPLT failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
		}
	}
	if cpu.ActiveTask().X < 0 {
		cpu.ActiveTask().IP = operand - 1
		// fmt.Printf("Jump to %#08x\n", operand)
	}
	return nil
}

func jmpgt(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return eris.Wrapf(err, "instruction JMPGT failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
		}
	}
	if cpu.ActiveTask().X > 0 {
		cpu.ActiveTask().IP = operand - 1
		// fmt.Printf("Jump to %#08x\n", operand)
	}
	return nil
}

func fetch(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return eris.Wrapf(err, "instruction FETCH failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
		}
	}
	if operand >= uint32(len(cpu.Ram)) {
		return eris.Wrapf(eris.New(fmt.Sprintf("available memory: %#08x", len(cpu.Ram))), "instruction FETCH failed at [%03d] %#08x: Invalid memory address %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP, operand)
	}
	value := cpu.Ram[operand]

	if err := cpu.push(value); err != nil {
		return eris.Wrapf(err, "instruction FETCH failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}
	return nil
}

func store(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return eris.Wrapf(err, "instruction STORE failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
		}
	}

	if operand >= uint32(len(cpu.Ram)) {
		return eris.Wrapf(eris.New(fmt.Sprintf("available memory: %#08x", len(cpu.Ram))), "instruction STORE failed at [%03d] %#08x: Invalid memory address %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP, operand)
	}

	value, err := cpu.pop()
	if err != nil {
		return eris.Wrapf(err, "instruction STORE failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
	}

	cpu.Ram[operand] = value
	return nil
}

func in(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func out(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func spawn(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return eris.Wrapf(err, "instruction SPAWN failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
		}
	}
	cpu.SpawnTask(operand)
	return nil
}

func yield(cpu *nova64Cpu, operand uint32) error {
	cpu.yield()
	return nil
}

func wait(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return eris.Wrapf(err, "instruction WAIT failed at [%03d] %#08x", cpu.ActiveTask().ID, cpu.ActiveTask().IP)
		}
	}
	cpu.ActiveTask().Waiting = operand
	cpu.yield()
	return nil
}

func kill(cpu *nova64Cpu, operand uint32) error {
	cpu.ActiveTask().killed = true
	return nil
}
