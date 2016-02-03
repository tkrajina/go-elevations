# SRTM parser for golang

go-elevations is a parser for "The Shuttle Radar Topography Mission" data.

It is based on the existing library for python [srtm.py](https://github.com/tkrajina/srtm.py)

## Usage

    package main

    import (
        "fmt"

        "github.com/tkrajina/go-elevations/geoelevations"
    )

    func main() {
        srtm, err := geoelevations.NewSrtm()
        if err != nil {
            panic(err.Error())
        }
        elevation, err := srtm.GetElevation(45.2775, 13.726111)
        if err != nil {
            panic(err.Error())
        }
        fmt.Println("Višnjan elevation is", elevation)
    }

## License

This library is licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)
