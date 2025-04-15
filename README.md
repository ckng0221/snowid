# SnowID - Distributed Unique ID Generator

SnowID is module that generates 64-bit unique ID in a distributed system based on Twitter Snowflake ID.
The generated ID is unique, sortable, and can be generated distributedly without single point of failure.
The ID is generated in binary string, and convert to decimal integer string

Structure of ID:
|Component|0|timestamp|datacenterID|machineID|sequenceNumber|
|---|---|---|---|---|---|
|Bit|1|41|5|5|12|

- 0: Placeholder, could indicate the sign, remain for future use.
- timestamp: unix timestamp in milliseconds. The default epoch is on 2025-01-01T00:00.000Z. 41 bit could have total of 2^41 = 2.2 Trillion milliseconds ~= 69.8 years. Max date could reach year 2095.
- dataCenterID: max 32 data center
- machineId: max 32 machine (nodes) in 1 data center
- sequenceNUmber: max 4096 sequence number per millisecond

Example generated ID:
(DataCenterID: 13, MachineID: 14, Epoch: DefaultEpoch)

- Binary: `0000001001011101101101111111111110100010010110101110000000000000`
- Decimal: `170494669478354944`

## Getting started

### Installation

```bash
go get github.com/ckng0221/snowid
```

### Quickstart

```go
package main

import (
	"fmt"

	"github.com/ckng0221/snowid"
)

func main() {
	dataCenterId := 1            // 0 to 31
	machineId := 1               // 0 to 31
	epoch := snowid.DefaultEpoch // Default epoch 2025-01-01T00:00.000Z
	s, err := snowid.NewSnowIdGenerator(dataCenterId, machineId, epoch)
	if err != nil {
		panic(err)
	}
	id1 := s.GenerateId()
	id2 := s.GenerateId()

	// ID 1
	fmt.Printf("ID: %s\n", id1.String())                // output, eg. 37866498659848192
	fmt.Printf("ID (Binary): %s\n", id1.StringBinary()) // output, eg. 0000000010000110100001110110000101000001100000100001000000000000
	fmt.Printf("Sequence: %d\n", id1.SequenceNumber)    // sequence 0
	// ID 2
	fmt.Printf("ID: %s\n", id2.String())                // outpuot, eg. 37866498659848193
	fmt.Printf("ID (Binary): %s\n", id2.StringBinary()) // output, eg. 0000000010000110100001110110000101000001100000100001000000000001
	fmt.Printf("Sequence: %d\n", id2.SequenceNumber)    // sequence 1

	// Parse ID
	id1Copy := id1.String()
	reverseParseId1, _ := snowid.ParseId(id1Copy, snowid.DefaultEpoch)
	fmt.Printf("ID: %s\n", reverseParseId1.String()) // same id as ID 1 after parsing, ie. 37866498659848192
}
```

## Contributing

We welcome your contribution!

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on submitting patches and the contribution workflow.
