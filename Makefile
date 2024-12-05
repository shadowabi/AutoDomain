BUILD_FLAG = -trimpath -ldflags "-s -w"
BUILD_DIR = bin

default: build

tidy:
	go mod tidy

build linux: tidy
	env CGO_ENABLED=0  GOOS=linux GOARCH=amd64 go build $(BUILD_FLAG) -o $(BUILD_DIR)/AutoDomain main/main.go

build windows: tidy
	env CGO_ENABLED=0  GOOS=windows GOARCH=amd64 go build $(BUILD_FLAG) -o $(BUILD_DIR)/AutoDomain.exe main/main.go

clean:
	rm -rf ./$(BUILD_DIR)
