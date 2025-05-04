package main

import (
	"fmt"

	"github.com/ckng0221/snowid"
)

func main() {
	dataCenterID := 1            // 0 to 31
	machineID := 1               // 0 to 31
	epoch := snowid.DefaultEpoch // Default epoch 2025-01-01T00:00.000Z
	s := snowid.NewSnowIDGenerator(dataCenterID, machineID, epoch)
	id1, err := s.GenerateID()
	if err != nil {
		panic(err)
	}
	id2, err := s.GenerateID()
	if err != nil {
		panic(err)
	}

	// ID 1
	fmt.Printf("ID (string): %s\n", id1.String())       // output, eg. "37866498659848192"
	fmt.Printf("ID (Binary): %s\n", id1.StringBinary()) // output, eg. "0000000010000110100001110110000101000001100000100001000000000000"
	fmt.Printf("ID (Decimal): %d\n", id1.Int64())       // output, eg. 37866498659848192
	fmt.Printf("Sequence: %d\n", id1.SequenceNumber)    // sequence 0
	// ID 2
	fmt.Printf("ID: %s\n", id2.String())                // outpuot, eg. "37866498659848193"
	fmt.Printf("ID (Binary): %s\n", id2.StringBinary()) // output, eg. "0000000010000110100001110110000101000001100000100001000000000001"
	fmt.Printf("ID (Decimal): %d\n", id2.Int64())       // output, eg. 37866498659848192
	fmt.Printf("Sequence: %d\n", id2.SequenceNumber)    // sequence 1

	// Parse ID
	id1Copy := id1.String()
	reverseParseID1, _ := snowid.ParseID(id1Copy, snowid.DefaultEpoch)
	fmt.Printf("ID: %s\n", reverseParseID1.String()) // same id as ID 1 after parsing, ie. "37866498659848192"
}
