package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/rotisserie/eris"
)

var opCodeMap = map[string]byte{
	"NOOP":  0x00,
	"PUSH":  0x01,
	"DROP":  0x02,
	"DUP":   0x03,
	"SWAP":  0x04,
	"ADD":   0x10,
	"SUB":   0x11,
	"MUL":   0x12,
	"DIV":   0x13,
	"MOD":   0x14,
	"AND":   0x20,
	"OR":    0x21,
	"XOR":   0x22,
	"NOT":   0x23,
	"CMP":   0x30,
	"JMP":   0x40,
	"CALL":  0x41,
	"RET":   0x42,
	"JMPZ":  0x43,
	"JMPLT": 0x44,
	"JMPGT": 0x45,
	"FETCH": 0x50,
	"STORE": 0x51,
	"IN":    0x60,
	"OUT":   0x61,
	"SPAWN": 0x70,
	"YIELD": 0x71,
	"WAIT":  0x72,
	"IRET":  0x80,
	"TRAP":  0x81,
	"KILL":  0xFF,
}

func assemble(path string) ([]byte, error) {
	// result := make([]byte, 0)
	var buffer bytes.Buffer
	writer := bufio.NewWriter(&buffer)
	labels := make(map[string]uint32)

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	counter := uint32(0)

	// collect labels
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, ":") {
			parts := strings.Fields(line)

			label := strings.TrimPrefix(parts[0], ":")
			label = strings.ToUpper(label)
			labels[label] = counter

			if len(parts) > 1 {
				counter++
			}

		} else if line != "" && !strings.HasPrefix(line, ";") {
			counter++
		}
	}

	// write labels to result
	labelCount := uint32(len(labels))
	if err = binary.Write(writer, binary.LittleEndian, labelCount); err != nil {
		return nil, err
	}
	for label, address := range labels {
		if err = binary.Write(writer, binary.LittleEndian, address); err != nil {
			return nil, err
		}
		_, err = writer.WriteString(label)
		if err != nil {
			return nil, err
		}
		// result = append(result, label...)
		writer.WriteByte(0)
		// result = append(result, 0)
	}

	// assemble
	file.Seek(0, 0)
	scanner = bufio.NewScanner(file)
	counter = 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		counter++
		if line == "" || strings.HasPrefix(line, ";") {
			continue
		}
		if strings.HasPrefix(line, ":") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}
			line = strings.Join(parts[1:], " ")
		}

		value, err := assembleLine(line, labels)
		if err != nil {
			return nil, eris.Wrapf(err, "error in line %d", counter)
		}
		if err = binary.Write(writer, binary.LittleEndian, value); err != nil {
			return nil, err
		}

	}

	writer.Flush()
	return buffer.Bytes(), nil
}

func assembleLine(line string, labels map[string]uint32) (uint32, error) {
	parts := strings.Fields(line)
	parts[0] = strings.ToUpper(parts[0])
	opCode, ok := opCodeMap[parts[0]]

	if !ok {
		return 0, eris.New(fmt.Sprintf("OpCode '%s' not found", parts[0]))
	}

	// no operand
	if len(parts) == 1 {
		return uint32(opCode), nil
	}

	// numeral operand
	arg := strings.ToUpper(parts[1])
	if v, err := strconv.ParseInt(arg, 0, 24); err == nil {
		return uint32(v)<<8 | uint32(opCode), nil
	}

	// check for size if numeral operand could not be parsed
	if _, err := strconv.ParseInt(arg, 0, 32); err == nil {
		return 0, eris.New(fmt.Sprintf("%s exceeds operand size", arg))
	}

	// label operand
	if v, ok := labels[arg]; ok {
		return v<<8 | uint32(opCode), nil
	}

	// TOP operand
	if arg == "TOP" {
		return 0xFFFFFF<<8 | uint32(opCode), nil
	}

	return 0, eris.New(fmt.Sprintf("unknown operand '%s'", arg))
}
