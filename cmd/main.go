package main

import (
	"fmt"
	"os"
	"utils/internal/cli"
)

func main() {
	if err := cli.HandleCommand(os.Stdout, os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "main error: %v\n", err)
		os.Exit(1)
	}
}

// func main() {
// 	opt, err := config.ParseFlags(os.Stdout, os.Args[1:])
// 	if err != nil {
// 		fmt.Fprintln(os.Stdout, err)
// 		os.Exit(1)
// 	}

// 	fmt.Println("Mode:", opt.Mode)
// 	fmt.Println("Env:", opt.Env)
// }
