package common

import (
	"io"
	"net/http"
)

func RemoteGET(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	var resBody []byte
	if err != nil {
		Logger.Error("http client: could not create request: ", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		Logger.Error("http client: could not send request: ", err)
	} else {
		resBody, err = io.ReadAll(res.Body)
		// resBody bytes to string
		if err != nil {
			Logger.Error("http client: could not read response: ", err)
		}
	}
	return resBody, err
}
