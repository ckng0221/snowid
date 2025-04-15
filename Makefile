test_all: test_unit test_build_basic test_build_server

test_unit:
	go test .

# Test examples
test_build_basic:
	cd examples/basic/
	go test -c -o test_build_basic ./...
	./test_build_basic
	rm test_build_basic

test_build_server:
	cd examples/server/
	go test -c -o test_build_server ./...
	./test_build_server
	rm test_build_server