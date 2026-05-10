BINARY    := pomodoro
BIN       := bin
SOUND_SRC := cmd/pomodoro/cli/sounds/done.wav

.PHONY: build run sound clean linux darwin-arm darwin-amd

build: $(SOUND_SRC)
	mkdir -p $(BIN)
	go build -ldflags="-s -w" -o $(BIN)/$(BINARY) ./cmd/pomodoro/

run: $(SOUND_SRC)
	go run ./cmd/pomodoro/

sound: $(SOUND_SRC)

$(SOUND_SRC):
	go run scripts/gen_sound.go

linux: $(SOUND_SRC)
	mkdir -p $(BIN)
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN)/$(BINARY)-linux-amd64 ./cmd/pomodoro/

darwin-arm: $(SOUND_SRC)
	mkdir -p $(BIN)
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o $(BIN)/$(BINARY)-darwin-arm64 ./cmd/pomodoro/

darwin-amd: $(SOUND_SRC)
	mkdir -p $(BIN)
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BIN)/$(BINARY)-darwin-amd64 ./cmd/pomodoro/

clean:
	rm -rf $(BIN)
	rm -f $(SOUND_SRC)
