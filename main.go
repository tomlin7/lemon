package main

import (
	"fmt"
	"lemon/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! Welcome to \x1b[48;5;226mLemon REPL\x1b[0m\n", user.Username)
	fmt.Println("Type in commands.")
	repl.Start(os.Stdin, os.Stdout)
}
