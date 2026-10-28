package main

import (
	"fmt"
	"github.com/Amir-Korouni/portman/internal/launcher"
	"github.com/Amir-Korouni/portman/internal/ports"
	"os"
	"flag"
	"errors"
	"context"
	"text/tabwriter"
)

const maxAttempts = 5

func printOccupied(from, to int) error {
	entries, err := ports.Occupied(context.Background(), from, to)
	if err != nil {
		return err
	}
	if len(entries) == 0 {
		fmt.Printf("no occupied ports in range %d-%d\n", from, to)
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, "PORT\tPID\tPROCESS\tADDRESS")
	for _, e := range entries {
		pid, name := "-", "-" 
		if e.PID > 0 {
			pid = fmt.Sprint(e.PID)
		}
		if e.Process != "" {
			name = e.Process
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", e.Port, pid, name, e.Address)
	}
	return w.Flush()
}

func main() {
	from := flag.Int("from", 3000, "first port of the range")
	to := flag.Int("to", 3999, "last port of the range")
	list := flag.Bool("list", false, "list occupied ports (with PID) in the range and exit")
	flag.Parse()

	if *list {
		if err := printOccupied(*from, *to); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: portman [--from N] [--to N] [--list] [--] <program> [args...]")
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
