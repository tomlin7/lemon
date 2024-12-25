package main

import (
	"fmt"
	"lemon/cli"
	"lemon/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	args := os.Args[1:]
	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			os.Exit(1)
		}

		defer file.Close()
		cli.Exec(file, os.Stdout)
	} else {

		fmt.Printf("Hello %s! Welcome to \x1b[48;5;226mLemon REPL\x1b[0m\n", user.Username)
		fmt.Println("Type in commands.")
		repl.Start(os.Stdin, os.Stdout)
	}
}
