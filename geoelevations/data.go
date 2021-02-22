package geoelevations

// Info (to be (se)serialized) about all the SRTM files and their URLs
type SrtmData struct {
	name, description, baseURL string
	files                      map[string]string
}

func (sd SrtmData) getFileURL(fn string) (string, bool) {
	path, found := sd.files[fn]
	if !found {
		return "", false
	}
	return sd.baseURL + path, true
}
