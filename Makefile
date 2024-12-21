say_hello:
	@echo "Hello. I'm a make file. I'm not sure why I'm here."

build:
	go build -o TimecodeTool ./cmd/TimecodeTool/main.go

test_not_valid:
	@./TimecodeTool 29.97 "00:07:00;00"

test_span:
	@./TimecodeTool 23.98 "01:00:00:00" "01:01:00:00"

test:
	@go test ./timecode -v
