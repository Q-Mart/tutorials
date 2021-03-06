package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"testing"

	"github.com/ReconfigureIO/sdaccel/xcl"
)

func BenchmarkKernel(world xcl.World, b *testing.B) {
	// Get our program
	program := world.Import("kernel_test")
	defer program.Release()

	// Get our kernel
	krnl := program.GetKernel("reconfigure_io_sdaccel_builder_stub_0_1")
	defer krnl.Release()

	// We need to create an input the size of B.N, so that the kernel
	// iterates B.N times
	input := make([]uint32, b.N)

	// seed it with random values, bound to 0 - 2**16
	for i, _ := range input {
		input[i] = uint32(uint16(rand.Uint32()))
	}

	// Create input buffer
	buff := world.Malloc(xcl.ReadOnly, uint(binary.Size(input)))
	defer buff.Free()

	// Create input buffer
	resp := make([]byte, b.N)
	outputBuff := world.Malloc(xcl.ReadWrite, uint(binary.Size(resp)))
	defer outputBuff.Free()

	// Write input buffer
	binary.Write(buff.Writer(), binary.LittleEndian, &input)

	// Clear output buffer
	binary.Write(outputBuff.Writer(), binary.LittleEndian, &resp)

	// Set args
	krnl.SetMemoryArg(0, buff)
	krnl.SetMemoryArg(1, outputBuff)
	krnl.SetArg(2, uint32(len(input)))

	// Reset the timer so that we only measure runtime of the kernel
	b.ResetTimer()
	krnl.Run(1, 1, 1)
}

func main() {
	// Create the world
	world := xcl.NewWorld()
	defer world.Release()

	// Create a function that the benchmarking machinery can call
	f := func(b *testing.B) {
		BenchmarkKernel(world, b)
	}

	// Benchmark it
	result := testing.Benchmark(f)

	log.Println()
	log.Println()
	log.Printf("FPGA runtime benchmark: ")
	log.Println()
	// Print the result
	fmt.Printf("%s\n", result.String())
}
