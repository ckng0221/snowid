# Distributed Unique ID Generator Server

An example server that generates 64-bit unique ID in a distributed system.
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

To run the server:

```bash
go run main.go
```

To test the APIs:

```bash
# Create Unique ID
curl http://localhost:8000/unique-ids -X POST

# Get Unique ID
# eg.
id="132343457927323648"
curl http://localhost:8000/unique-ids/$id
```
