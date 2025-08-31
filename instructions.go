package main

import (
	"errors"
	"fmt"

	comp "cmp"
)

const TOP uint32 = 0xFFFFFF

func noop(cpu *nova64Cpu, operand uint32) error {
	return nil
}

func push(cpu *nova64Cpu, operand uint32) error {
	if err := cpu.push(operand); err != nil {
		return fmt.Errorf("instruction PUSH failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func drop(cpu *nova64Cpu, operand uint32) error {
	if err := cpu.drop(); err != nil {
		return fmt.Errorf("instruction DROP failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func dup(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return fmt.Errorf("instruction DUP failed at %#08x: %w", cpu.ActiveTask.IP, err)
		}
	}
	if err := cpu.dup(operand); err != nil {
		return fmt.Errorf("instruction DUP failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func swap(cpu *nova64Cpu, operand uint32) error {
	if err := cpu.swap(); err != nil {
		return fmt.Errorf("instruction SWAP failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func add(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction ADD failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	b, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction ADD failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	err = cpu.push(a + b)
	if err != nil {
		return fmt.Errorf("instruction ADD failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func sub(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction SUB failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	b, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction SUB failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	err = cpu.push(a - b)
	if err != nil {
		return fmt.Errorf("instruction SUB failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func mul(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction MUL failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	b, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction MUL failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	err = cpu.push(a * b)
	if err != nil {
		return fmt.Errorf("instruction MUL failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func div(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction DIV failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	b, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction DIV failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	err = cpu.push(a / b)
	if err != nil {
		return fmt.Errorf("instruction DIV failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func mod(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction MOD failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	b, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction MOD failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	err = cpu.push(a % b)
	if err != nil {
		return fmt.Errorf("instruction MOD failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func and(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction AND failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	b, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction AND failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	err = cpu.push(a & b)
	if err != nil {
		return fmt.Errorf("instruction AND failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func or(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction OR failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	b, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction OR failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	err = cpu.push(a | b)
	if err != nil {
		return fmt.Errorf("instruction OR failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func xor(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction XOR failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	b, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction XOR failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	err = cpu.push(a ^ b)
	if err != nil {
		return fmt.Errorf("instruction XOR failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func not(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction XOR failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	err = cpu.push(^a)
	if err != nil {
		return fmt.Errorf("instruction XOR failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	return nil
}

func cmp(cpu *nova64Cpu, operand uint32) error {
	a, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction CMP failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}
	b, err := cpu.pop()
	if err != nil {
		return fmt.Errorf("instruction CMP failed at %#08x: %w", cpu.ActiveTask.IP, err)
	}

	cpu.ActiveTask.X = int32(comp.Compare(a, b))
	return nil
}

func jmp(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return fmt.Errorf("instruction JMP failed at %#08x: %w", cpu.ActiveTask.IP, err)
		}
	}
	cpu.ActiveTask.IP = operand - 1
	// fmt.Printf("Jump to %#08x\n", operand)
	return nil
}

func call(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func ret(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func jmpz(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return fmt.Errorf("instruction JMPZ failed at %#08x: %w", cpu.ActiveTask.IP, err)
		}
	}
	if cpu.ActiveTask.X == 0 {
		cpu.ActiveTask.IP = operand - 1
		// fmt.Printf("Jump to %#08x\n", operand)
	}
	return nil
}

func jmplt(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return fmt.Errorf("instruction JMPLT failed at %#08x: %w", cpu.ActiveTask.IP, err)
		}
	}
	if cpu.ActiveTask.X < 0 {
		cpu.ActiveTask.IP = operand - 1
		// fmt.Printf("Jump to %#08x\n", operand)
	}
	return nil
}

func jmpgt(cpu *nova64Cpu, operand uint32) error {
	if operand == TOP {
		var err error = nil
		operand, err = cpu.pop()
		if err != nil {
			return fmt.Errorf("instruction JMPGT failed at %#08x: %w", cpu.ActiveTask.IP, err)
		}
	}
	if cpu.ActiveTask.X > 0 {
		cpu.ActiveTask.IP = operand - 1
		// fmt.Printf("Jump to %#08x\n", operand)
	}
	return nil
}

func fetch(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func store(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func in(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func out(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func spawn(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func yield(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func wait(cpu *nova64Cpu, operand uint32) error {
	panic(errors.New("NotImplemented"))
}

func halt(cpu *nova64Cpu, operand uint32) error {
	cpu.Halted = true
	return nil
}
