package geoelevations

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

func prepareClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}, nil
}

func downloadSrtmURL(url, username, password string) ([]byte, error) {
	client, err := prepareClient()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return loadResp(resp)
	}

	if resp.StatusCode == 302 {
		loc := resp.Header["Location"][0]
		fmt.Println("Redirecting to:", loc)

		authReq, err := http.NewRequest(http.MethodGet, loc, nil)
		if err != nil {
			return nil, err
		}

		authReq.SetBasicAuth(username, password)

		authResp, err := client.Do(authReq)
		if err != nil {
			return nil, err
		}

		fmt.Println("Auth resp status:", authResp.Status)

		codeLoc := authResp.Header["Location"][0]
		fmt.Println("Redirect to:", codeLoc)

		codeReq, err := http.NewRequest(http.MethodGet, codeLoc, nil)
		if err != nil {
			return nil, err
		}

		codeResp, err := client.Do(codeReq)
		if err != nil {
			return nil, err
		}

		fmt.Println("Code resp status: ", codeResp.Status)

		lastRedirectLoc := codeResp.Header["Location"][0]
		fmt.Println("Last redirect:", lastRedirectLoc)
	}

	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return loadResp(resp)
}

func loadResp(resp *http.Response) ([]byte, error) {
	byts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return byts, nil
}
