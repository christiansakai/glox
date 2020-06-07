package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/christiansakai/glox/scanner"
)

var hadError = false

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
	} else if len(os.Args) == 2 {
		err := runFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
		}
	} else {
		runPrompt()
	}
}

func runFile(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error reading file: %w", err)
	}

	run(string(file))
	if hadError {
		os.Exit(65)
	}

	return nil
}

func runPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		str, _ := reader.ReadString('\n')
		run(str)
		hadError = false
	}
}

func run(source string) {
	scn := scanner.New(source)
	tokens := scn.ScanTokens(reportError)

	for _, token := range tokens {
		fmt.Println(token)
	}
}

func reportError(line int, message string) {
	report(line, "", message)
}

func report(line int, where, message string) {
	fmt.Println("[line %d] Error %s: %s", line, where, message)
	hadError = true
}
