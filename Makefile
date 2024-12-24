say_hello:
	@echo "Hello. I'm a make file. I'm not sure why I'm here."

build:
	go build -o TimecodeTool ./cmd/TimecodeTool/main.go

build_wasm:
	GOOS=js GOARCH=wasm go build -o dist/timecodetool.wasm ./cmd/wasm/main.go
	cp cmd/wasm/wasm_exec.js dist/
	cp cmd/wasm/index.html dist/

build_wasm_tinygo:
	#this is not working yet.
	tinygo build -o dist/timecodetool_tiny.wasm -target wasm ./cmd/wasm/main.go
	cp cmd/wasm/wasm_exec.js dist/
	cp cmd/wasm/index.html dist/


test_not_valid:
	@./TimecodeTool 29.97 "00:07:00;00"

test_span:
	@./TimecodeTool 23.98 "01:00:00:00" "01:01:00:00"

test:
	@go test ./timecode -v
