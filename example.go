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
