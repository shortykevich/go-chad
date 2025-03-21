build-exec:
	go mod tidy && go build -o app

run-local-server:
	./app -mode=server -host=localhost -port=8554

run-client:
	./app -mode=client -host=localhost -port=8554
