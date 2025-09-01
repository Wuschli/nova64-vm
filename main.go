package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	cpu := NewCpu(4 * 1024 * 1024)
	fmt.Printf("RAM size: %d words, %d bytes\n", len(cpu.Ram), len(cpu.Ram)*4)

	start := time.Now()
	data, err := assemble("test.nasm")
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	fmt.Printf("% x\n", data)
	fmt.Printf("Assembling took %s\n", elapsed)

	if err = cpu.Load(data); err != nil {
		panic(err)
	}

	ticker := time.NewTicker(time.Microsecond)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				cpu.Tick()
				if cpu.Halted {
					cancel()
				}
			}
		}
	}()

	<-ctx.Done()
	ticker.Stop()
}
