package main

import (
	"fmt"
	"os"
	"github.com/Amir-Korouni/portman/internal/launcher"
)

func main() {
	fmt.Println("Hello, Go!")
	if len(os.Args) < 2 {
		fmt.Println("usage: go run , <program> [args...]")
		os.Exit(1)
	}
	if err := launcher.RunProgram(os.Args[1], os.Args[2:]...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
