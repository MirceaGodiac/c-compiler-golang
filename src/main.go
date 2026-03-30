package main

import "fmt"

func main() {
	sourceCode := `
		int main_var = 42;
		main_var++;
		if (main_var == 43) {
			return 0;
		}

		for (int i = 0; i < 10; i++) {
			return i;
		}
		return 0;
	`

	tokens := Tokenize(sourceCode)

	for _, t := range tokens {
		fmt.Println(t)
	}
}
