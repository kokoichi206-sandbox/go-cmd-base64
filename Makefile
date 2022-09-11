BINARY_NAME=base64

## build: Build binary
build:
	env CGO_ENABLED=0 go build -ldflags="-s -w" -o ${BINARY_NAME} .

## clean: runs go clean and deletes binaries
clean:
	@go clean
	@rm ${BINARY_NAME}
	@echo "Cleaned!"

## test: runs all tests
test:
	go test -v ./...
