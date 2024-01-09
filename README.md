# eflag - Environment Variable Flag Overrides

`eflag` is designed to override flag values using environment variables when they are not set explicitly.

Flag value precedence (from greatest to least):

1. explicitly set flag
2. environment variable
3. flag default value

## Installation

```shell
go get github.com/zachcheung/eflag
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/zachcheung/eflag"
)

func main() {
	var (
		host string
		port int
	)

	eflag.Var(&host, "host", "localhost", "host", "")
	eflag.Var(&port, "port", 8000, "port", "-")

	eflag.SetPrefix("myapp") // Optionally set prefix for environment variables
	eflag.Parse()

	fmt.Println("host:", host) // Environment variable "MYAPP_HOST" will overrides default "localhost" if it's set
	fmt.Println("port:", port) // Ignores environment variable "MYAPP_PORT"
}
```

## License

[MIT](LICENSE)
