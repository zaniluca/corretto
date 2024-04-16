# Corretto

Corretto (Italian for "Free from errors") is a simple and essential schema validation package for go structs. It is designed with code readability in mind.

## Getting Started

run the following Go command to install the corretto package:

```bash
go get github.com/zaniluca/corretto
```

## Usage

```go
package main

import (
	"fmt"
	"github.com/zaniluca/corretto"
)

type User struct {
	Name string
	Age  int
}

func main() {
	user := User{
		Name: "John Doe",
		Age:  17,
	}

	schema := corretto.Schema{
		"Name": corretto.Field().Required(),
		"Age":  corretto.Field().Required().Min(18),
	}

	err := corretto.Validate(user)
	if err != nil {
		fmt.Println(err) // ValidationError: Age must be greater than or equal to 18
	}
}
```

## Documentation
