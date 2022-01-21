# go-stp

Go-stp provides bindings for the [Simple Theorem Prover](https://stp.github.io/) on Linux.
Most operations are supported.

## Installation

The bindings make use of the [Simple Theorem Prover](https://stp.github.io/) which should be setup first.
After STP has been setup, just clone the repo in your $GOPATH/src folder.

```bash
git clone https://github.com/arckt/go-stp.git
```

Note that for the bindings to work, the stp dynamic library is assumed to be in /usr/local/lib/.
If the library is unaccessible run the following.

```bash
sudo ldconfig /path/to/lib/library
```

## Usage

The bindings could be used purely as a c wrapper which is done through cgo or with the convenient Python-like API which is used in the Python STP bindings.

An example of a golang program that uses the Python-like API.

```go
package main

import "fmt"
import g"arckt/go-stp"

func main() {
  h := g.Init() //Initialize solver
  defer h.Destroy() //Destroy solver at end of main()

  x := h.Bitvec("x", 32) //Add symbollic variable of bitwidth 32 of name "x"

  h.add("x + 3 == 5") //Add constraint for the solver to evaluate
  res := h.Solve(x) //Tell solver to evaluate value of symbollic variable x to fulfill constraints

  fmt.Println(res[0]) //Print symbollic variable of x if it is solvable
}
```

# Contributing

Anyone is free to contribute or suggest any changes to this repository if they are interested.

# License

[MIT](https://choosealicense.com/licenses/mit/)
