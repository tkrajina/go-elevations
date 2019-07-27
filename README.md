# SRTM parser for golang

go-elevations is a parser for "The Shuttle Radar Topography Mission" data.

It is based on the existing library for python [srtm.py](https://github.com/tkrajina/srtm.py)

## Usage

```golang
package main

import (
	"fmt"
	"net/http"

	"github.com/tkrajina/go-elevations/geoelevations"
)

func main() {
	srtm, err := geoelevations.NewSrtm(http.DefaultClient)
	if err != nil {
		panic(err.Error())
	}
	elevation, err := srtm.GetElevation(http.DefaultClient, 45.2775, 13.726111)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Višnjan elevation is", elevation)
}
```

go-elevations is a parser for "The Shuttle Radar Topography Mission" data.

It is based on the existing library for python [srtm.py](https://github.com/tkrajina/srtm.py)

## Usage

http.DefaultClient

## License

This library is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

## License

This library is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
