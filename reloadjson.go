package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/tkrajina/geoelevations/geoelevations"
)

func main() {
	srtmData, err := geoelevations.LoadSrtmData()

	if err != nil {
		log.Panic("Error reloading json:", err.Error())
		return
	}

	srtmDataJson, err := json.MarshalIndent(srtmData, "", "\t")
	if err != nil {
		log.Panic("Error marshalling srtmData:", err.Error())
		return
	}

	fileName := "geoelevations/urls.json"
	f, err := os.Create(fileName)
	if err != nil {
		log.Panic("Error writing json to ", fileName, ":", err.Error())
		return
	}

	f.Write(srtmDataJson)

	log.Print("Written ", len(srtmDataJson), " bytes to ", fileName)
}
