package main

import (
	"bufio"
	"fmt"
	"lox-go/scanner"
	"os"
)

func main() {
	args := os.Args[1:]

	if len(args) > 1 {
		fmt.Println("Usage: jlox [script]")
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

func runFile(filepath string) {

	filebyte, err := os.ReadFile(filepath)

	if err != nil {
		panic(err)
	}

	run(string(filebyte))
}

func runPrompt() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for scanner.Scan() == true {
		line := scanner.Text()

		run(line)
		fmt.Print("> ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input: ", err)
	}
}

func run(line string) {
	scanner := scanner.NewScanner(line)
	scanner.ScanTokens()

	for _, token := range scanner.Tokens {
		fmt.Println(token)
	}
}

