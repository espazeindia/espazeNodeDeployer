.PHONY: run build test clean install-deps

run:
	go run cmd/server/main.go

build:
	go build -o bin/espazeNodeDeployer cmd/server/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/

install-deps:
	go mod download
	go mod tidy

dev:
	air -c .air.toml

docker-build:
	docker build -t espaze-node-deployer:latest .

docker-run:
	docker run -p 8080:8080 espaze-node-deployer:latest

