run_server_dev:
	compileDaemon -command="./example"

run_server:
	go run main.go

test_id:
	curl -X POST http://localhost:8000/snowids

test_id_with_load:
	autocannon -c 10 -a 2000 -m "POST" http://localhost:8000/snowids

test_get_id:
	curl http://localhost:8000/unique-ids/170135204736458752
