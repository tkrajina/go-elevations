package geoelevations

import (
	"fmt"
	"testing"
)

func checkSrtmFileName(t *testing.T, latitude, longitude float64, expectedFileName string) {
	srtm := NewSrtm()
	fileName := srtm.getSrtmFileName(latitude, longitude)
	if fileName != expectedFileName {
		t.Error(fmt.Sprintf("SRTM FILE for (%v, %v) should be %s but is %s", latitude, longitude, expectedFileName, fileName))
	}
}

func TestFindSrtmFileName(t *testing.T) {
	checkSrtmFileName(t, 45, 13, "N45E013.hgt")
	checkSrtmFileName(t, 45.1, 13, "N45E013.hgt")
	checkSrtmFileName(t, 44.9, 13, "N44E013.hgt")
	checkSrtmFileName(t, 45, 13.1, "N45E013.hgt")
	checkSrtmFileName(t, 45, 12.9, "N45E012.hgt")
	checkSrtmFileName(t, 25, -80, "N25W080.hgt")
	checkSrtmFileName(t, 25, -80.1, "N25W081.hgt")
	checkSrtmFileName(t, 25, -79.9, "N25W080.hgt")
	checkSrtmFileName(t, 25.1, -80, "N25W080.hgt")
	checkSrtmFileName(t, -32, 152, "S32E152.hgt")

	// This file don't exists but the get_file_name is expected to return the supposed file:
	checkSrtmFileName(t, 0, 0, "N00E000.hgt")
}
