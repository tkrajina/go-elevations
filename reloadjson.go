package main

import (
	"github.com/tkrajina/geoelevations/geoelevations"
	"log"
	"os"
)

func main() {
	json, err := geoelevations.GetSrtmFilesUrls()
	if err != nil {
		log.Panic("Error reloading json:", err.Error())
		return
	}

	fileName := "geoelevations/urls.json"
	f, err := os.Create(fileName)
	if err != nil {
		log.Panic("Error writing json to ", fileName, ":", err.Error())
		return
	}

	f.Write(json)

	log.Print("Written ", len(json), " bytes to ", fileName)
}
