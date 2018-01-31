package geoelevations

import (
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkSrtmFileName(t *testing.T, latitude, longitude float64, expectedFileName string, expectedSrtmLatitude, expectedSrtmLongitude float64) {
	srtm, _ := NewSrtm()
	fileName, srtmLatitude, srtmLongitude := srtm.getSrtmFileNameAndCoordinates(latitude, longitude)
	log.Printf("Checking %s", fileName)
	if fileName != expectedFileName {
		t.Error(fmt.Sprintf("SRTM FILE for (%v, %v) should be %s but is %s", latitude, longitude, expectedFileName, fileName))
	}
	if srtmLatitude != expectedSrtmLatitude {
		t.Errorf("srtmLatitude != expectedSrtmLatitude ... %f != %f", srtmLatitude, expectedSrtmLatitude)
	}
	if srtmLongitude != expectedSrtmLongitude {
		t.Errorf("srtmLongitude != expectedSrtmLongitude ... %f != %f", srtmLongitude, expectedSrtmLongitude)
	}
}

func TestFindSrtmFileName(t *testing.T) {
	checkSrtmFileName(t, 45, 13, "N45E013", 45, 13)
	checkSrtmFileName(t, 45.1, 13, "N45E013", 45, 13)
	checkSrtmFileName(t, 44.9, 13, "N44E013", 44, 13)
	checkSrtmFileName(t, 45, 13.1, "N45E013", 45, 13)
	checkSrtmFileName(t, 45, 12.9, "N45E012", 45, 12)
	checkSrtmFileName(t, 25, -80, "N25W080", 25, -80)
	checkSrtmFileName(t, 25, -80.1, "N25W081", 25, -81)
	checkSrtmFileName(t, 25, -79.9, "N25W080", 25, -80)
	checkSrtmFileName(t, 25.1, -80, "N25W080", 25, -80)
	checkSrtmFileName(t, -32, 152, "S32E152", -32, 152)

	// This file don't exists but the get_file_name is expected to return the supposed file:
	checkSrtmFileName(t, 0, 0, "N00E000", 0, 0)
}

func checkElevation(t *testing.T, latitude, longitude, expectedElevation float64) {
	srtm, _ := NewSrtm()
	elevation, err := srtm.GetElevation(http.DefaultClient, latitude, longitude)
	if err != nil {
		t.Errorf("Valid coordinates but error getting elevation:%s", err.Error())
		return
	}
	if elevation != expectedElevation {
		t.Errorf("Invalid elevation for (%f, %f): %f, but should be %f", latitude, longitude, elevation, expectedElevation)
	}
}

func TestGetElevation(t *testing.T) {
	checkElevation(t, 45.2775, 13.726111, 246)
	checkElevation(t, -26.4, 146.25, 301)
	checkElevation(t, -12.1, -77.016667, 133)
	checkElevation(t, 40.75, -111.883333, 1298)
}

func TestSrtm1AndSrtm3ForUSA(t *testing.T) {
	srtm, _ := NewSrtm()
	srtmFileName, _, _ := srtm.getSrtmFileNameAndCoordinates(40.75, -111.883333)
	storage, err := NewLocalFileSrtmStorage("")
	assert.Nil(t, err)
	srtmData, err := newSrtmData(storage)
	assert.Nil(t, err)
	srtm1url := srtmData.GetSrtm1Url(srtmFileName)
	srtm3url := srtmData.GetSrtm3Url(srtmFileName)
	if len(srtm1url.Url) == 0 || len(srtm3url.Url) == 0 {
		t.Error("USA should have both srtm1 and srtm3 urls")
	}
}

func TestSrtm1AndSrtm3ForEurope(t *testing.T) {
	srtm, _ := NewSrtm()
	srtmFileName, _, _ := srtm.getSrtmFileNameAndCoordinates(45.2775, 13.726111)
	storage, err := NewLocalFileSrtmStorage("")
	assert.Nil(t, err)
	srtmData, err := newSrtmData(storage)
	assert.Nil(t, err)
	srtm1url := srtmData.GetSrtm1Url(srtmFileName)
	srtm3url := srtmData.GetSrtm3Url(srtmFileName)
	if !(srtm3url != nil && len(srtm3url.Url) > 0 && srtm1url == nil) {
		t.Error("Europe should have both srtm1 and srtm3 urls")
	}
}
