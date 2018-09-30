package main

import (
	// import the entire framework (including bundled verilog)
	_ "github.com/ReconfigureIO/sdaccel"

	// Use the SMI protocol package for interacting with shared memory
	"github.com/ReconfigureIO/sdaccel/smi"
)

func Mult(a uint32) uint32 {
	return a * 2
}

func Top(
	// Specify inputs and outputs to and from the FPGA. Tell the FPGA where to find data in shared memory, what data type
	// to expect or pass single integers directly to the FPGA by sending them to the control register - see examples

	inputData uintptr,
	outputData uintptr,
	length uint32,

	// Set up ports for interacting with the shared memory, here we have 2 SMI ports which can be used to read or write
	readReq chan<- smi.Flit64,
	readResp <-chan smi.Flit64,

	writeReq chan<- smi.Flit64,
	writeResp <-chan smi.Flit64) {

	// Do whatever needs doing with the data from the host
	inputChan := make(chan uint32)
	go smi.ReadBurstUInt32(readReq, readResp, inputData, smi.DefaultOptions, length, inputChan)

	transformedChan := make(chan uint32)
	go func() {
		for i := length; i > 0; i-- {
			num := <-inputChan
			transformedChan <- Mult(num)
		}
	}()

	// Write the result to the location in shared memory as requested by the host
	smi.WriteBurstUInt32(writeReq, writeResp, outputData, smi.DefaultOptions, length, transformedChan)
}
