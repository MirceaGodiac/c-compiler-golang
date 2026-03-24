# C Compiler in Go

A C compiler written from scratch in Go. Currently implements the **lexer** (lexical analysis) stage.

## What's Implemented

### Lexer

The lexer takes C source code and breaks it into a stream of tokens by scanning character-by-character.

**Supported tokens:**

| Category | Tokens |
|---|---|
| Keywords | `int`, `return` |
| Identifiers | variable and function names (e.g. `main_var`) |
| Integer literals | decimal numbers (e.g. `42`) |
| Operators | `=`, `==`, `+`, `++` |
| Punctuation | `;` |

Unrecognized characters (e.g. `(`, `)`, `{`, `}`, `if`) are flagged as illegal tokens.

## Project Structure

```
src/
  token.go              - Token types, Token struct, keyword lookup
  lexical-analyser.go   - Lexer: character scanning and tokenization
  main.go               - Entry point with a hardcoded C snippet for testing
```

## Build & Run

```bash
cd src
go build -o c-lexer.exe
./c-lexer.exe
```
