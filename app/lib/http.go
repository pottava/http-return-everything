package lib

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"
)

const retryCount = 1

// HTTPGet send a HTTP request
func HTTPGet(ctx context.Context, client *http.Client, endpoint string, headers *http.Header) ([]byte, error) {
	var err error
	for i := 0; i < retryCount; i++ {
		resp, e := httpGetOnce(ctx, client, http.MethodGet, endpoint, headers)
		if e == nil {
			return resp, nil
		}
		err = e
		time.Sleep(100 * time.Millisecond)
	}
	return nil, err
}

// HTTPPut send a HTTP request
func HTTPPut(ctx context.Context, client *http.Client, endpoint string, headers *http.Header) ([]byte, error) {
	var err error
	for i := 0; i < retryCount; i++ {
		resp, e := httpGetOnce(ctx, client, http.MethodPut, endpoint, headers)
		if e == nil {
			return resp, nil
		}
		err = e
		time.Sleep(100 * time.Millisecond)
	}
	return nil, err
}

func httpGetOnce(ctx context.Context, client *http.Client, method, endpoint string, headers *http.Header) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, endpoint, nil)
	if err != nil {
		return nil, err
	}
	if headers != nil {
		req.Header = *headers
	}
	resp, err := client.Do(req)
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
