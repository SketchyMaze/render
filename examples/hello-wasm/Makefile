.PHONY: build
build:
	GOOS=js GOARCH=wasm go build -v -o hello.wasm main_wasm.go

.PHONY: serve
serve: build
	go run server.go

.PHONY: clean
clean:
	rm -f hello.wasm
