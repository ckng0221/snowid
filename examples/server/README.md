# Distributed Unique ID Generator Server

An example server that uses `snowid` to generate the unique ID.

## Getting started

To run the server:

```bash

go run main.go
```

To test the APIs:

```bash
# Create Unique ID
curl http://localhost:8000/ids -X POST

# Get Unique ID
# eg.
id="132343457927323648"
curl http://localhost:8000/ids/$id
```
