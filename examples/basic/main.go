package main

import (
	"fmt"

	"github.com/ckng0221/snowid"
)

func main() {
	dataCenterId := 1
	machineId := 1
	g, err := snowid.NewSnowIdGenerator(dataCenterId, machineId, snowid.DefaultEpoch)
	if err != nil {
		panic(err)
	}
	id1 := g.GenerateId()
	id2 := g.GenerateId()
	id4 := g.GenerateId()
	fmt.Println(id4)

	// ID 1
	fmt.Printf("ID: %s\n", id1.String())
	fmt.Printf("ID (Binary): %s\n", id1.StringBinary())
	fmt.Printf("Sequence: %d\n", id1.SequenceNumber)
	// ID 2
	fmt.Printf("ID: %s\n", id2.String())
	fmt.Printf("ID (Binary): %s\n", id2.StringBinary())
	fmt.Printf("Sequence: %d\n", id2.SequenceNumber)

	// Parse ID
	id1Copy := id1.String()
	fmt.Println(id1Copy)
	id3, _ := snowid.ParseId(id1Copy, snowid.DefaultEpoch)
	fmt.Printf("ID: %s\n", id3.String())
}
