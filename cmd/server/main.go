package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
)

func run(
	ctx context.Context,
	args []string,
	getenv func(string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	fmt.Println("Hello, World")
	return nil
}

func main() {
	ctx := context.Background()
	if err := run(
		ctx,
		nil,
		nil,
		nil,
		os.Stdout,
		os.Stderr,
	); err != nil {
		log.Fatalf("%s", err)
	}
}
