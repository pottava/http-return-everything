package lib

import (
	"io/ioutil"
	"net/http"
	"time"
)

func HTTPGet(client *http.Client, endpoint string) ([]byte, error) {
	var err error
	for i := 0; i < 3; i++ {
		resp, e := httpGetOnce(client, endpoint)
		if e == nil {
			return resp, nil
		}
		err = e
		time.Sleep(100 * time.Millisecond)
	}
	return nil, err
}

func httpGetOnce(client *http.Client, endpoint string) ([]byte, error) {
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
