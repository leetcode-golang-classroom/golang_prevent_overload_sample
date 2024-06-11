.PHONY=build

build-rate-limiter:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/rate-limiter rate-limiter/cmd/main.go

run-rate-limiter: build-rate-limiter
	@./bin/rate-limiter

build-load-shedding:
	@CGO_ENABLED=0 GOOS=linux go build -o bin/load-shedding load-shedding/cmd/main.go

run-load-shedding: build-load-shedding
	@./bin/load-shedding
	
coverage:
	@go test -v -cover ./...

test:
	@go test -v ./...
