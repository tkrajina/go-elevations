package geoelevations

import (
	"fmt"
	"log"
	"testing"
)

func checkSrtmFileName(t *testing.T, latitude, longitude float64, expectedFileName string, expectedSrtmLatitude, expectedSrtmLongitude float64) {
	srtm := NewSrtm()
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

func TestGetElevation(t *testing.T) {
	srtm := NewSrtm()
	elevation, err := srtm.GetElevation(45.2775, 13.726111)
	if err != nil {
		t.Errorf("Valid coordinates but error getting elevation:%s", err.Error())
		return
	}
	if elevation != 246 {
		t.Errorf("Invalid elevation for Višnjan:%v", elevation)
	}
}
