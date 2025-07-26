package main

import (
	"fmt"
	"os"

	"com.loop.anonx3247/src/lexer"
)

func main() {
	// Check if a file path is provided as command line argument
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file-path>")
		fmt.Println("Please provide a file path to lex")
		os.Exit(1)
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

func printTokens(tokens []lexer.Token) {
	for i, token := range tokens {
		if token.Type == lexer.EOF {
			fmt.Printf("%d: EOF\n", i+1)
		} else if token.Type == lexer.NEWLINE {
			fmt.Printf("%d: NEWLINE\n", i+1)
		} else {
			fmt.Printf("%d: %d ('%s')\n", i+1, int(token.Type), token.Value.String())
		}
	}
}
