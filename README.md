# godress
Parses street addresses.

## Installing

Install the Go way:

```sh
go get -u github.com/ecarter202/godress
```

## Using

`````go
package main

import (
	"fmt"

	ap "github.com/ecarter202/godress"
)

func main() {
	s := "5723 NE Golang Ave Gopherville, UT 39232"

	a, _ := ap.Parse(s)

	fmt.Printf("%+v\n", *a)

  // Output:
  // {HouseNumber:5723 StreetDirection:NE StreetName:Golang StreetType:Ave Unit: City:Gopherville State:UT Zip:39232}
}
`````
