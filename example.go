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
