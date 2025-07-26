package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"com.loop.anonx3247/lexer"
	"com.loop.anonx3247/parser"
)

func main() {
	// Check if a file path is provided as command line argument
	if len(os.Args) < 2 {
		// Enter REPL mode
		runREPL()
		return
	}

	// Get the file path from command line arguments
	filePath := os.Args[1]

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file '%s': %v\n", filePath, err)
		os.Exit(1)
	}

	// Convert content to string
	source := string(content)

	// Create a new lexer with the file content
	lex := lexer.NewLexer(source)

	// Tokenize the content
	tokens, err := lex.Tokenize()
	if err != nil {
		fmt.Printf("Error tokenizing file: %v\n", err)
		printTokens(tokens)
		os.Exit(1)
	}

	// Display the tokens
	fmt.Printf("Tokenizing file: %s\n", filePath)
	fmt.Printf("Found %d tokens:\n\n", len(tokens))

	printTokens(tokens)
}

func runREPL() {
	fmt.Println("Loop Language REPL")
	fmt.Println("Type 'exit' or 'quit' to exit, or press Ctrl+C")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("loop> ")

		if !scanner.Scan() {
			// EOF or error
			break
		}

		input := scanner.Text()
		input = strings.TrimSpace(input)

		// Check for exit commands
		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// Skip empty lines
		if input == "" {
			continue
		}

		p, err := parser.NewParser(input)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		expr, err := p.ParseExpr()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		val, err := expr.Eval()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Println(val)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %v\n", err)
	}
}

func printTokens(tokens []lexer.Token) {
	for i, token := range tokens {
		if token.Type == lexer.EOF {
			fmt.Printf("%d: EOF  ", i+1)
		} else if token.Type == lexer.NEWLINE {
			fmt.Printf("%d: NEWLINE ", i+1)
		} else {
			fmt.Printf("%d: %d ('%s')  ", i+1, int(token.Type), token.Value.String())
		}
	}
	fmt.Println()
}
