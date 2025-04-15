run_server:
	go run main.go

run_server_dev:
	compileDaemon -command="./example"

test_id:
	curl -X POST http://localhost:8000/ids

test_id_with_load:
	autocannon -c 10 -a 2000 -m "POST" http://localhost:8000/ids

get_id:
	curl http://localhost:8000/ids/170494669478354944
