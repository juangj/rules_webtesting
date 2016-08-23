package httphelper

import (
	"context"
	"io/ioutil"
	"net/http"
)

var client = &http.Client{}

// Forward forwards r to host and writes the response from host to w.
func Forward(ctx context.Context, host string, w http.ResponseWriter, r *http.Request) error {
	url := "http://" + host + r.URL.Path

	// Construct request based on Method, URL Path, and Body from r
	request, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return err
	}
	request = request.WithContext(ctx)

	// Copy headers from r to request
	request.Header = r.Header

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Copy response headers from resp to w
	for k, vs := range resp.Header {
		w.Header().Del(k)
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}

	// Copy status code from resp to w
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
	return nil
}

func Get(ctx context.Context, url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	request = request.WithContext(ctx)
	return client.Do(request)
}