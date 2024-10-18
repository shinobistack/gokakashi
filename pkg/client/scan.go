package client

import (
	"net/http"
)

type ScanOpts struct {
	ApiHeaders map[string]string
	Severity   string
}

func Scan(imageName, apiHost string, opts *ScanOpts) error {
	// TODO: Complete this function after completing
	// https://github.com/gokakashi/goKakashi/issues/22

	req, err := http.NewRequest(http.MethodPost, apiHost, nil)
	if err != nil {
		return err
	}

	for k, v := range opts.ApiHeaders {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
