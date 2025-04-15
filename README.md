# Distributed Unique ID Generator

A POC of gnerating 64-bit unique ID in a distributed system.
The generated ID is unique, sortable, and can be generated distributedly without single point of failure.
The ID is generated in binary string, and convert to decimal integer string

Structure of ID:
|Component|0|timestamp|datacenterID|machineID|sequenceNumber|
|---|---|---|---|---|---|
|Bit|1|41|5|5|12|

- 0: Placeholder, could indicate the sign, remain for future use.
- timestamp: unix timestamp in milliseconds. Have custom epoch of 1704067200000 (starting from 2024-01-01T00:00.000Z). 41 bit could have total of 2^41 = 2.2 Trillion milliseconds ~= 69.8 years. Max date could reach year 2093.
- dataCenterID: max 32 data center
- machineId: max 32 machine (nodes) in 1 data center
- sequenceNUmber: max 4096 sequence number per millisecond

Example generated ID:

- Binary: `0110010001000010101011010010000010000001000110111110000000000110`
- Decimal: `7224527107372343302`

## Getting started

### Prerequisites

Gin requires [Go](https://go.dev/) version [1.24](https://go.dev/doc/devel/release#go1.24.0) or above.

### Getting module

With [Go's module support](https://go.dev/wiki/Modules#how-to-use-modules), `go [build|run|test]` automatically fetches the necessary dependencies when you add the import in your code:

```sh
import "github.com/ckng0221/distributed-unique-id-generator"
```

### Running module

A basic example:

```go
package main

import (
  "github.com/ckng0221/distributed-unique-id-generator"
)

func main() {
  r := gin.Default()
  r.GET("/ping", func(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
      "message": "pong",
    })
  })
  r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
```
