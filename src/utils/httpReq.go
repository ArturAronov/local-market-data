package utils

import (
	"fmt"
	"io"
	"net/http"
)

func HttpReq(email string, url string) ([]byte, error) {
	client := &http.Client{}

	req, reqErr := http.NewRequest("GET", url, nil)
	if reqErr != nil {
		return nil, reqErr
	}

	req.Header.Set("User-Agent", email)

	res, resErr := client.Do(req)
	if resErr != nil {
		return nil, resErr
	}

	defer res.Body.Close()

	body, bodyErr := io.ReadAll(res.Body)
	if bodyErr != nil {
		return nil, bodyErr
	}

	if res.StatusCode >= 300 {
		return nil, fmt.Errorf("Status code: %d, body: %s\n", res.StatusCode, string(body))
	}

	return body, nil
}
