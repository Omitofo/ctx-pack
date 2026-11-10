package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("ctx-pack %s\n", version)
		os.Exit(0)
	}

	fmt.Println("ctx-pack — package folders into clean context for LLMs")
	fmt.Println("Status: early development")
	fmt.Println()
	fmt.Println("Usage will be available soon.")
}
