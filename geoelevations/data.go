package geoelevations

type SrtmUrl struct {
	File string
	Url  string
}

type SrtmData struct {
	Srtm1 []SrtmUrl
	Srtm3 []SrtmUrl
}
