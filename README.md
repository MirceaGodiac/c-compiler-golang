# C Compiler in Go

A C compiler written from scratch in Go. Currently implements the **lexer** (lexical analysis) stage.

## What's Implemented

### Lexer

The lexer takes C source code and breaks it into a stream of tokens using a **table-driven deterministic finite automaton (DFA)**.

Instead of hand-written `switch` logic, the lexer defines:

- **Character classes** that group input bytes into equivalence classes (letter, digit, `=`, `+`, `-`, `;`, other)
- **DFA states** representing progress through a token (e.g. `StateStart → StateAssign → StateEq`)
- **A transition table** (`transition[state][charClass] → nextState`) initialized at startup
- **An accept table** (`acceptToken[state] → TokenType`) mapping accepting states to token types
- **A driver loop** that runs the DFA with longest-match semantics and rewind support

Adding a new token requires only new states and table entries — no new control flow.

**Supported tokens:**

| Category | Tokens |
|---|---|
| Keywords | `int`, `float`, `string`, `var`, `double`, `if`, `return` |
| Identifiers | variable and function names (e.g. `main_var`) |
| Integer literals | decimal numbers (e.g. `42`) |
| Operators | `=`, `==`, `+`, `++` |
| Punctuation | `;` |

Unrecognized characters (e.g. `(`, `)`, `{`, `}`) are flagged as illegal tokens.

## Project Structure

```
src/
  token.go              - Token types, Token struct, keyword lookup
  lexical-analyser.go   - Table-driven DFA lexer: states, transitions, and driver loop
  main.go               - Entry point with a hardcoded C snippet for testing
```

## Build & Run

```bash
cd src
go build -o c-lexer.exe
./c-lexer.exe
```
