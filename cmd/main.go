package main

import (
	"apps/config"
	"fmt"
	"os"
)

func main() {
	opt, err := config.ParseFlags(os.Stdout, os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}

	fmt.Println("Mode:", opt.Mode)
	fmt.Println("Env:", opt.Env)
}
