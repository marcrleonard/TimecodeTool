//go:build js && wasm

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"syscall/js"

	"TimecodeTool/handlers"
)

func ProcessInput(this js.Value, args []js.Value) interface{} {
	if len(args) == 0 {
		return js.ValueOf("Error: no input provided")
	}
	input := args[0].String()
	fmt.Println("Input received:", input) // Log the input to the console

	resp := handlers.ValidateTimecode(input, 29.97)

	var buf bytes.Buffer

	// Pass the pointer to the buffer to the encoder
	if err := json.NewEncoder(&buf).Encode(resp); err != nil {
		return js.ValueOf("Error encoding JSON")
	}

	return js.ValueOf(buf.String())
}

func main() {
	// Expose the `processInput` function to JavaScript
	js.Global().Set("processInput", js.FuncOf(ProcessInput))

	// Keep the Go program running
	select {}
}
