package geoelevations

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

type httpClient struct {
	http.Client
	username, password string
}

func newHTTPClient(username, password string) (*httpClient, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &httpClient{
		Client: http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
			Jar: jar,
		},
		username: username,
		password: password,
	}, nil
}

func (c *httpClient) downloadSrtmURL(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return c.loadResp(resp)
	}

	if resp.StatusCode == 302 {
		loc := resp.Header["Location"][0]
		fmt.Println("Redirecting to:", loc)

		authReq, err := http.NewRequest(http.MethodGet, loc, nil)
		if err != nil {
			return nil, err
		}

		authReq.SetBasicAuth(c.username, c.password)

		authResp, err := c.Do(authReq)
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

		codeResp, err := c.Do(codeReq)
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
	resp, err = c.Do(req)
	if err != nil {
		return nil, err
	}

	return c.loadResp(resp)
}

func (c *httpClient) loadResp(resp *http.Response) ([]byte, error) {
	byts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return byts, nil
}
