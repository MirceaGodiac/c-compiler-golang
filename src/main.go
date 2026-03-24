package main

import "fmt"

func main() {
	sourceCode := `
		int main_var = 42;
		main_var++;
		if (main_var == 43) {
			return 0;
		}
	`
	// Note: Our lexer doesn't handle 'if', '(', ')', '{', '}' yet.
	// Watch how it flags them as Illegal.

	tokens := Tokenize(sourceCode)

	for _, t := range tokens {
		fmt.Println(t)
	}
}
