package main

import (
	"errors"
	"fmt"

	"github.com/rotisserie/eris"

	comp "cmp"
)

const TOP uint32 = 0xFFFFFF

func noop(cpu *nova64Cpu) error {
	return nil
}

func push(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "PUSH")
	}
	if err := cpu.push(operand); err != nil {
		return cpu.wrapError(err, "PUSH")
	}
	return nil
}

func drop(cpu *nova64Cpu) error {
	if err := cpu.drop(); err != nil {
		return cpu.wrapError(err, "PUSH")
	}
	return nil
}

func dup(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "DUP")
	}
	if err := cpu.dup(operand); err != nil {
		return cpu.wrapError(err, "DUP")
	}
	return nil
}

func swap(cpu *nova64Cpu) error {
	if err := cpu.swap(); err != nil {
		return cpu.wrapError(err, "SAWP")
	}
	return nil
}

func add(cpu *nova64Cpu) error {
	b, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "ADD")
	}
	a, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "ADD")
	}
	err = cpu.push(uint32(int32(a) + int32(b)))
	if err != nil {
		return cpu.wrapError(err, "ADD")
	}
	return nil
}

func sub(cpu *nova64Cpu) error {
	b, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "SUB")
	}
	a, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "SUB")
	}
	err = cpu.push(uint32(int32(a) - int32(b)))
	if err != nil {
		return cpu.wrapError(err, "SUB")
	}
	return nil
}

func mul(cpu *nova64Cpu) error {
	b, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "MUL")
	}
	a, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "MUL")
	}
	err = cpu.push(uint32(int32(a) * int32(b)))
	if err != nil {
		return cpu.wrapError(err, "MUL")
	}
	return nil
}

func div(cpu *nova64Cpu) error {
	b, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "DIV")
	}
	a, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "DIV")
	}
	err = cpu.push(uint32(int32(a) / int32(b)))
	if err != nil {
		return cpu.wrapError(err, "DIV")
	}
	return nil
}

func mod(cpu *nova64Cpu) error {
	b, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "MOD")
	}
	a, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "MOD")
	}
	err = cpu.push(uint32(int32(a) % int32(b)))
	if err != nil {
		return cpu.wrapError(err, "MOD")
	}
	return nil
}

func and(cpu *nova64Cpu) error {
	b, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "AND")
	}
	a, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "AND")
	}
	err = cpu.push(a & b)
	if err != nil {
		return cpu.wrapError(err, "AND")
	}
	return nil
}

func or(cpu *nova64Cpu) error {
	b, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "OR")
	}
	a, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "OR")
	}
	err = cpu.push(a | b)
	if err != nil {
		return cpu.wrapError(err, "OR")
	}
	return nil
}

func xor(cpu *nova64Cpu) error {
	b, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "XOR")
	}
	a, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "XOR")
	}
	err = cpu.push(a ^ b)
	if err != nil {
		return cpu.wrapError(err, "XOR")
	}
	return nil
}

func not(cpu *nova64Cpu) error {
	a, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "NOT")
	}
	err = cpu.push(^a)
	if err != nil {
		return cpu.wrapError(err, "NOT")
	}
	return nil
}

func cmp(cpu *nova64Cpu) error {
	b, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "CMP")
	}
	a, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "CMP")
	}

	cpu.ActiveTask().X = int32(comp.Compare(int32(a), int32(b)))
	return nil
}

func jmp(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "JMP")
	}
	cpu.ActiveTask().IP = operand - 1
	// fmt.Printf("Jump to %#08x\n", operand)
	return nil
}

func call(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "CALL")
	}

	err = cpu.push(cpu.ActiveTask().IP)
	if err != nil {
		return cpu.wrapError(err, "CALL")
	}

	cpu.ActiveTask().IP = operand - 1
	return nil
}

func ret(cpu *nova64Cpu) error {
	value, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "RET")
	}
	cpu.ActiveTask().IP = value - 1
	return nil
}

func jmpz(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "JMPZ")
	}
	if cpu.ActiveTask().X == 0 {
		cpu.ActiveTask().IP = operand - 1
		// fmt.Printf("Jump to %#08x\n", operand)
	}
	return nil
}

func jmplt(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "JMPLT")
	}
	if cpu.ActiveTask().X < 0 {
		cpu.ActiveTask().IP = operand - 1
		// fmt.Printf("Jump to %#08x\n", operand)
	}
	return nil
}

func jmpgt(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "JMPGT")
	}
	if cpu.ActiveTask().X > 0 {
		cpu.ActiveTask().IP = operand - 1
		// fmt.Printf("Jump to %#08x\n", operand)
	}
	return nil
}

func fetch(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "FETCH")
	}
	if operand >= uint32(len(cpu.Ram)) {
		return cpu.wrapError(eris.New(fmt.Sprintf("available memory: %#08x", len(cpu.Ram))), "FETCH")
	}
	value := cpu.Ram[operand]

	if err := cpu.push(value); err != nil {
		return cpu.wrapError(err, "FETCH")
	}
	return nil
}

func store(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "STORE")
	}

	if operand >= uint32(len(cpu.Ram)) {
		return cpu.wrapError(eris.New(fmt.Sprintf("available memory: %#08x", len(cpu.Ram))), "STORE")
	}

	value, err := cpu.pop()
	if err != nil {
		return cpu.wrapError(err, "STORE")
	}

	cpu.Ram[operand] = value
	return nil
}

func in(cpu *nova64Cpu) error {
	panic(errors.New("NotImplemented"))
}

func out(cpu *nova64Cpu) error {
	panic(errors.New("NotImplemented"))
}

func spawn(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "SPAWN")
	}
	cpu.SpawnTask(operand)
	return nil
}

func yield(cpu *nova64Cpu) error {
	cpu.yield()
	return nil
}

func wait(cpu *nova64Cpu) error {
	operand, err := cpu.fetchOperand()
	if err != nil {
		return cpu.wrapError(err, "WAIT")
	}
	cpu.ActiveTask().Waiting = operand
	cpu.yield()
	return nil
}

func iret(cpu *nova64Cpu) error {
	panic(errors.New("NotImplemented"))
}

func trap(cpu *nova64Cpu) error {
	panic(errors.New("NotImplemented"))
}

func kill(cpu *nova64Cpu) error {
	cpu.ActiveTask().killed = true
	return nil
}
