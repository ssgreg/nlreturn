# nlreturn

[![build status](https://github.com/ssgreg/nlreturn/actions/workflows/main.yaml/badge.svg?branch=master)]()
[![Go Report Status](https://goreportcard.com/badge/github.com/ssgreg/nlreturn)](https://goreportcard.com/report/github.com/ssgreg/nlreturn)
[![Coverage Status](https://coveralls.io/repos/github/ssgreg/nlreturn/badge.svg?branch=master&service=github)](https://coveralls.io/github/ssgreg/nlreturn?branch=master)

Linter requires a new line before return and branch statements except when the return is alone inside a statement group (such as an if statement) to increase code clarity.

# Example

Examples of incorrect code:

```go
func foo() int {
    a := 0
    _ = a
    return a
}

func bar() int {
    a := 0
    if a == 0 {
        _ = a
        return
    }
    return a
}
```

Examples of correct code:

```go
func foo() int {
    a := 0
    _ = a

    return a
}

func bar() int {
    a := 0
    if a == 0 {
        _ = a

        return
    }

    return a
}
```
