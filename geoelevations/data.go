package geoelevations

// Info (to be (se)serialized) about all the SRTM files and their URLs
type SrtmData struct {
	BaseURL string            `json:"baseUrl"`
	Files   map[string]string `json:"files"`
}

func (sd SrtmData) getFileURL(fn string) (string, bool) {
	path, found := urls.Files[fn]
	if !found {
		return "", false
	}
	return urls.BaseURL + path, true
}
