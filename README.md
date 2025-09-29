# idl-parser

OMG IDL Parser written in go inspired by [gomme](https://github.com/oleiade/gomme)

An OMG IDL (Interface Definition Language) parser written in Go, inspired by gomme.

## Features
* Parses OMG IDL syntax into structured Go types
* Supports common IDL constructs:
  * Modules
  * Basic Types
  * Type references
  * Arrays/Sequence
  * Annotations
* Simple API with Parse() function
* Comprehensive test coverage

## Basic Type Transform

| IDL Base Type        | Common Base Type | Description             |
|----------------------|------------------|-------------------------|
| `octet`              | `uint8`          | 8-bit unsigned integer  |
| `short`              | `int16`          | 16-bit signed integer   |
| `unsigned short`     | `uint16`         | 16-bit unsigned integer |
| `long`               | `int32`          | 32-bit signed integer   |
| `unsigned long`      | `uint32`         | 32-bit unsigned integer |
| `long long`          | `int64`          | 64-bit signed integer   |
| `unsigned long long` | `uint64`         | 64-bit unsigned integer |
| `float`              | `float32`        | 32-bit floating point   |
| `double`             | `float64`        | 64-bit floating point   |
| `boolean`            | `bool`           | Boolean value           |
| `string`             | `string`         | Variable-length string  |

```txt
module m {
     struct c {
        octet id1;
        short id2;
        unsigned short id3;
     }
}
```

### Dynamic String

```txt
module m {
     struct c {
        string name;
     }
}
```

### Fixed Length String

```txt
module m {
     struct c {
        string<10> name;
     }
}
```

### Array 

Array is the fixed list of given element

```txt
module m {
     struct c {
        octet[10] idList;
     }
}
```

### Sequence

Array is the dynamic list of given element

```txt
module m {
     struct c {
        sequence<octet> idList;
     }
}
```

### TypeRef

```txt
module m {
    struct c1 {
        octet id1;
        short id2;
        unsigned short id3;
    }
    
    struct c2 {
        c1      refc1;
        short   id2;
    }
}
```

## Installation

```bash
go get github.com/Yisaer/idlparser
```

## Quick Start

```go
package main

import (
	"fmt"
	
	"github.com/Yisaer/idlparser/ast"
)

func main() {
	input := `module example {
		struct Point {
			long x;
			long y;
		};
	}`

	result := ast.Parse(input)
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
		@format(a="b",key=123) idbits id;
	};
}
```

See [ast_test.go](./ast/ast_test.go) for more parsing examples.

## License

This project is licensed under the terms of the MIT license. See [LICENSE](./LICENSE) for details.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.
