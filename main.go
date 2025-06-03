package main

import (
	"fmt"
	"io"
	"os"

	"github.com/yisaer/idl-parser/ast"
)

func main() {
	code, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	result := ast.Parse(string(code))
	if result.Err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing module: %v\n", result.Err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", result.Output)
}
