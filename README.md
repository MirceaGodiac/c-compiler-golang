# C Compiler in Go

A C compiler written from scratch in Go. Currently implements the **lexer** (lexical analysis) stage.

## Linter Setup

This project uses [golangci-lint](https://golangci-lint.run) and runs it automatically on every `git commit` via a pre-commit hook.

### Install golangci-lint

```sh
cd src && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

Then add the Go binary directory to your PATH (add this to `~/.bashrc` or `~/.zshrc`):

```sh
export PATH="$PATH:$HOME/go/bin"
```

### Run manually

```sh
cd src
golangci-lint run --config ../.golang-CI.yml ./...
```

The pre-commit hook runs this automatically and blocks the commit if any issues are found.

---

## What's Implemented

### Lexer

The lexer takes C source code and breaks it into a stream of tokens using a **table-driven deterministic finite automaton (DFA)**.

Instead of hand-written `switch` logic, the lexer defines:

- **Character classes** that group input bytes into equivalence classes
- **DFA states** representing progress through a token
- **A transition table** (`transition[state][charClass] → nextState`) initialized at startup
- **An accept table** (`acceptToken[state] → TokenType`) mapping accepting states to token types
- **A driver loop** that runs the DFA with longest-match semantics and rewind support

Adding a new token requires only new states and table entries — no new control flow.

---

## How the DFA Works

### 1. Character Classes (`token.go` / `lexical-analyser.go`)

Rather than branching on every possible byte value, input characters are first mapped into equivalence classes:

| Class | Matches |
|---|---|
| `ClassLetter` | `a-z`, `A-Z`, `_` |
| `ClassDigit` | `0-9` |
| `ClassEquals` | `=` |
| `ClassPlus` | `+` |
| `ClassMinus` | `-` |
| `ClassSemiColon` | `;` |
| `ClassStar` | `*` |
| `ClassSlash` | `/` |
| `ClassPercent` | `%` |
| `ClassLessThan` | `<` |
| `ClassGreaterThan` | `>` |
| `ClassExclamation` | `!` |
| `ClassAmpersand` | `&` |
| `ClassPipe` | `\|` |
| `ClassCaret` | `^` |
| `ClassTilde` | `~` |
| `ClassLParen` | `(` |
| `ClassRParen` | `)` |
| `ClassLBrace` | `{` |
| `ClassRBrace` | `}` |
| `ClassLBracket` | `[` |
| `ClassRBracket` | `]` |
| `ClassOther` | anything else |

This keeps the transition table small: `stateCount × classCount` cells instead of `stateCount × 256`.

### 2. DFA States

Each state represents a point in recognizing a token:

| State | Meaning |
|---|---|
| `StateStart` | No characters consumed yet |
| `StateAssign` | Seen `=` |
| `StateEq` | Seen `==` |
| `StatePlus` | Seen `+` |
| `StateMinus` | Seen `-` |
| `StateInc` | Seen `++` |
| `StateDec` | Seen `--` |
| `StateSemiColon` | Seen `;` |
| `StateIdent` | Inside an identifier `[A-Za-z_][A-Za-z0-9_]*` |
| `StateInt` | Inside an integer literal `[0-9]+` |
| `StateStar` | Seen `*` |
| `StateSlash` | Seen `/` |
| `StatePercent` | Seen `%` |
| `StateLessThan` | Seen `<` |
| `StateLessEq` | Seen `<=` |
| `StateLeftShift` | Seen `<<` |
| `StateGreaterThan` | Seen `>` |
| `StateGreaterEq` | Seen `>=` |
| `StateRightShift` | Seen `>>` |
| `StateNot` | Seen `!` |
| `StateNotEq` | Seen `!=` |
| `StateAmpersand` | Seen `&` |
| `StateAnd` | Seen `&&` |
| `StatePipe` | Seen `\|` |
| `StateOr` | Seen `\|\|` |
| `StateCaret` | Seen `^` |
| `StateTilde` | Seen `~` |
| `StateLParen` | Seen `(` |
| `StateRParen` | Seen `)` |
| `StateLBrace` | Seen `{` |
| `StateRBrace` | Seen `}` |
| `StateLBracket` | Seen `[` |
| `StateRBracket` | Seen `]` |

### 3. Transition Table

`transition[state][charClass]` gives the next state, or `NoTransition (-1)` if no valid move exists. The full table (all `NoTransition` cells omitted):

| From State | On Class | To State |
|---|---|---|
| `StateStart` | `ClassLetter` | `StateIdent` |
| `StateStart` | `ClassDigit` | `StateInt` |
| `StateStart` | `ClassEquals` | `StateAssign` |
| `StateStart` | `ClassPlus` | `StatePlus` |
| `StateStart` | `ClassMinus` | `StateMinus` |
| `StateStart` | `ClassSemiColon` | `StateSemiColon` |
| `StateStart` | `ClassStar` | `StateStar` |
| `StateStart` | `ClassSlash` | `StateSlash` |
| `StateStart` | `ClassPercent` | `StatePercent` |
| `StateStart` | `ClassLessThan` | `StateLessThan` |
| `StateStart` | `ClassGreaterThan` | `StateGreaterThan` |
| `StateStart` | `ClassExclamation` | `StateNot` |
| `StateStart` | `ClassAmpersand` | `StateAmpersand` |
| `StateStart` | `ClassPipe` | `StatePipe` |
| `StateStart` | `ClassCaret` | `StateCaret` |
| `StateStart` | `ClassTilde` | `StateTilde` |
| `StateStart` | `ClassLParen` | `StateLParen` |
| `StateStart` | `ClassRParen` | `StateRParen` |
| `StateStart` | `ClassLBrace` | `StateLBrace` |
| `StateStart` | `ClassRBrace` | `StateRBrace` |
| `StateStart` | `ClassLBracket` | `StateLBracket` |
| `StateStart` | `ClassRBracket` | `StateRBracket` |
| `StateAssign` | `ClassEquals` | `StateEq` |
| `StatePlus` | `ClassPlus` | `StateInc` |
| `StateMinus` | `ClassMinus` | `StateDec` |
| `StateLessThan` | `ClassEquals` | `StateLessEq` |
| `StateLessThan` | `ClassLessThan` | `StateLeftShift` |
| `StateGreaterThan` | `ClassEquals` | `StateGreaterEq` |
| `StateGreaterThan` | `ClassGreaterThan` | `StateRightShift` |
| `StateNot` | `ClassEquals` | `StateNotEq` |
| `StateAmpersand` | `ClassAmpersand` | `StateAnd` |
| `StatePipe` | `ClassPipe` | `StateOr` |
| `StateIdent` | `ClassLetter` | `StateIdent` |
| `StateIdent` | `ClassDigit` | `StateIdent` |
| `StateInt` | `ClassDigit` | `StateInt` |

### 4. Accept Table

When the DFA enters an accepting state, that state maps to the token it produces:

| Accepting State | Token Emitted |
|---|---|
| `StateAssign` | `TokenAssign` (`=`) |
| `StateEq` | `TokenEq` (`==`) |
| `StatePlus` | `TokenPlus` (`+`) |
| `StateMinus` | `TokenMinus` (`-`) |
| `StateInc` | `TokenInc` (`++`) |
| `StateDec` | `TokenDec` (`--`) |
| `StateSemiColon` | `TokenSemiColon` (`;`) |
| `StateIdent` | `TokenIdent` (or a keyword — see below) |
| `StateInt` | `TokenInt` |
| `StateStar` | `TokenStar` (`*`) |
| `StateSlash` | `TokenSlash` (`/`) |
| `StatePercent` | `TokenPercent` (`%`) |
| `StateLessThan` | `TokenLessThan` (`<`) |
| `StateLessEq` | `TokenLessEq` (`<=`) |
| `StateLeftShift` | `TokenLeftShift` (`<<`) |
| `StateGreaterThan` | `TokenGreaterThan` (`>`) |
| `StateGreaterEq` | `TokenGreaterEq` (`>=`) |
| `StateRightShift` | `TokenRightShift` (`>>`) |
| `StateNot` | `TokenExclamation` (`!`) |
| `StateNotEq` | `TokenNotEq` (`!=`) |
| `StateAmpersand` | `TokenAmpersand` (`&`) |
| `StateAnd` | `TokenAnd` (`&&`) |
| `StatePipe` | `TokenPipe` (`\|`) |
| `StateOr` | `TokenOR` (`\|\|`) |
| `StateCaret` | `TokenCaret` (`^`) |
| `StateTilde` | `TokenTilde` (`~`) |
| `StateLParen` | `TokenLParen` (`(`) |
| `StateRParen` | `TokenRParen` (`)`) |
| `StateLBrace` | `TokenLBrace` (`{`) |
| `StateRBrace` | `TokenRBrace` (`}`) |
| `StateLBracket` | `TokenLBracket` (`[`) |
| `StateRBracket` | `TokenRBracket` (`]`) |

`StateStart` and `NoTransition` are **not** accepting — hitting them without a prior accept means an illegal character.

### 5. Longest-Match Driver

The driver loop in `nextToken()` runs the DFA and tracks the last accepting state seen:

```
lastAcceptPos  = -1
lastAcceptType = TokenIllegal

loop:
    class := classify(currentChar)
    next  := transition[state][class]
    if next == NoTransition → break
    state = next
    if acceptToken[state] != TokenIllegal:
        record lastAcceptPos, lastAcceptType
    advance one character

if lastAcceptPos == -1 → emit TokenIllegal, advance one char
else:
    rewind to lastAcceptPos+1      ← so next call starts after the token
    lexeme = input[startPos : lastAcceptPos+1]
    emit token
```

This implements **maximal munch**: the DFA always consumes as many characters as possible before emitting a token. For example, `==` is scanned as a single `TokenEq`, not two `TokenAssign` tokens.

### 6. Keyword Resolution

When the DFA accepts an identifier (`TokenIdent`), the lexeme is looked up in the `keywords` map (`token.go`). If it matches a C keyword, the token type is replaced with the keyword's type (e.g. `TokenIntKw`, `TokenReturn`). All 33 standard C keywords are supported:

`auto` `break` `case` `char` `const` `continue` `default` `do` `double` `else` `enum` `extern` `float` `for` `goto` `if` `inline` `int` `long` `register` `restrict` `return` `short` `signed` `sizeof` `static` `struct` `switch` `typedef` `union` `unsigned` `void` `volatile` `while`

---

## Supported Tokens

| Category | Tokens |
|---|---|
| Keywords | All 33 C standard keywords (see above) |
| Identifiers | variable and function names (e.g. `main_var`) |
| Integer literals | decimal numbers (e.g. `42`) |
| Arithmetic operators | `=` `==` `+` `++` `-` `--` `*` `/` `%` |
| Comparison operators | `<` `<=` `>` `>=` `!=` |
| Shift operators | `<<` `>>` |
| Logical operators | `!` `&&` `\|\|` |
| Bitwise operators | `&` `\|` `^` `~` |
| Punctuation | `;` `(` `)` `{` `}` `[` `]` |

Unrecognized characters are flagged as `TokenIllegal`.

---

## Project Structure

```
src/
  token.go              - TokenType constants, Token struct, keywords map, LookupIdent
  lexical-analyser.go   - DFA states, character classes, transition/accept tables, driver loop
  main.go               - Entry point with a hardcoded C snippet for testing
```

## Build & Run

```bash
cd src
go build -o c-lexer.exe
./c-lexer.exe
```

## Linter Checks

```bash
    cd src // or wherever the code you want to idiot proof is...
    golangci-lint run
```