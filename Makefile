BINARY := pomodoro

.PHONY: build run clean linux darwin

build:
	go build -ldflags="-s -w" -o $(BINARY) ./cmd/pomodoro/

run:
	go run ./cmd/pomodoro/

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY)-linux-amd64 ./cmd/pomodoro/

darwin-arm:
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BINARY)-darwin-arm64 ./cmd/pomodoro/

darwin-amd:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BINARY)-darwin-amd64 ./cmd/pomodoro/

clean:
	rm -f $(BINARY) $(BINARY)-*
