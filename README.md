# idlparser

OMG IDL Parser written in go inspired by [gomme](https://github.com/oleiade/gomme)

An OMG IDL (Interface Definition Language) parser written in Go, inspired by gomme.

## Features
* Parses OMG IDL syntax into structured Go types
* Supports common IDL constructs:
  * Modules
  * Structs
  * Bitsets
  * Bitfields
  * Octet
  * Short / Unsigned Short
  * Long / Unsigned Long / Long Long / Unsigned Long Long
  * Sequence
  * Type references
  * Annotations
* Simple API with Parse() function
* Comprehensive test coverage

## Installation

```bash
go get github.com/Yisaer/idlparser
```

## Quick Start

```go
package main

import (
	"fmt"
	"github.com/Yisaer/idlparser"
)

func main() {
	input := `
	module example {
		struct Point {
			long x;
			long y;
		};
	}`

	result := idlparser.Parse(input)
	fmt.Printf("Parsed module: %+v\n", result.Output)
}
```

## Example

The parser can handle complex IDL definitions:

```bash
module spi {
	bitset idbits {
		bitfield<4> bid;  // 4 bits for bid
	};

	struct CANFrame {
		@format octet header;
		@format(a=b) idbits id;
	};
}
```

See [ast_test.go](./ast/ast_test.go) for more parsing examples.

## License

This project is licensed under the terms of the MIT license. See [LICENSE](./LICENSE) for details.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
