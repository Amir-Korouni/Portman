package main

import (
	"fmt"
	"github.com/Amir-Korouni/portman/internal/launcher"
	"github.com/Amir-Korouni/portman/internal/ports"
	"os"
	"flag"
	"errors"
)

const maxAttempts = 5

func main() {
	from := flag.Int("from", 3000, "first port of the range")
	to := flag.Int("to", 3999, "last port of the range")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: portman [--from N] [--to N] [--] <program> [args...]")
		flag.PrintDefaults()
		os.Exit(1)
	}


	start := *from
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if start > *to { 
			fmt.Fprintln(os.Stderr, ports.ErrNoFreePort)
			os.Exit(1)
		}

		port, err := ports.FirstFree(start, *to)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println("using port", port)

		err = launcher.RunProgram(args[0], port, args[1:]...)
		if err == nil {
			return
		}

		if errors.Is(err, launcher.ErrStartupFailed) {
			fmt.Fprintf(os.Stderr, "attempt %d on port %d failed: %v\n", attempt, port, err)
			start = port + 1 
			continue
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "giving up after %d attempts\n", maxAttempts)
	os.Exit(1)
}
