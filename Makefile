test_all: test_unit test_build_basic test_build_server

test_unit:
	go test ./...

# Test examples
test_build_basic:
	go build -o ./example-basic ./examples/basic/main.go
	rm example-basic

test_build_server:
	cd examples/server && \
	go build -o ./example-server ./main.go
	rm examples/server/example-server

# Run Examples
run_example_basic:
	go run examples/basic/main.go
run_example_server:
	go run examples/server/main.go