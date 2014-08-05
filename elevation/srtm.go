package geoelevations

import (
    "fmt"
    "math"
)

const (
    SRTM_BASE_URL = "http://dds.cr.usgs.gov/srtm"
    SRTM1_URL     = "/version2_1/SRTM1/"
    SRTM3_URL     = "/version2_1/SRTM3/"
)

func GetElevation(latitude, longitude float64) float64 {
    return 0
}

func getSrtmFileName(latitude, longitude float64) string {
    northSouth := 'S'
    if latitude >= 0 {
        northSouth ='N'
    }

    eastWest := 'W'
    if longitude >= 0 {
        eastWest = 'E'
    }

    latPart := int(math.Abs(math.Floor(latitude)))
    lonPart := int(math.Abs(math.Floor(longitude)))

    return fmt.Sprintf("%s%02d%s%03d.hgt", string(northSouth), latPart, string(eastWest), lonPart)
}
