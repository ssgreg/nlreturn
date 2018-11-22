# nlreturn

Linter requires a new line before return and branch statements to increase code clarity, except when the return is alone inside a statement group (such as an if statement).

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
