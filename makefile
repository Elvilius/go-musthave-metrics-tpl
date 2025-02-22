.PHONY: lint

run-agent: 
	go build -o agent cmd/agent/main.go && ./agent &

run-server: 
	go build -o server cmd/server/main.go && ./server &

lint:
	rm -f lint
	go build -o lint cmd/staticlint/main.go && ./lint ./...	