package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ipfs/ipfs-cluster/api"
)

type responseDecoder func(d *json.Decoder) error

func (c *defaultClient) do(
	method, path string,
	headers map[string]string,
	body io.Reader,
	obj interface{},
) error {

	resp, err := c.doRequest(method, path, headers, body)
	if err != nil {
		return &api.Error{Code: 0, Message: err.Error()}
	}
	return c.handleResponse(resp, obj)
}

func (c *defaultClient) doStream(
	method, path string,
	headers map[string]string,
	body io.Reader,
	outHandler responseDecoder,
) error {

	resp, err := c.doRequest(method, path, headers, body)
	if err != nil {
		return &api.Error{Code: 0, Message: err.Error()}
	}
	return c.handleStreamResponse(resp, outHandler)
}

func (c *defaultClient) doRequest(
	method, path string,
	headers map[string]string,
	body io.Reader,
) (*http.Response, error) {

	urlpath := c.net + "://" + c.hostname + "/" + strings.TrimPrefix(path, "/")
	logger.Debugf("%s: %s", method, urlpath)

	r, err := http.NewRequest(method, urlpath, body)
	if err != nil {
		return nil, err
	}
	if c.config.DisableKeepAlives {
		r.Close = true
	}

	if c.config.Username != "" {
		r.SetBasicAuth(c.config.Username, c.config.Password)
	}

	if headers != nil {
		for k, v := range headers {
			r.Header.Set(k, v)
		}
	}

	if body != nil {
		r.ContentLength = -1 // this lets go use "chunked".
	}

	return c.client.Do(r)
}
func (c *defaultClient) handleResponse(resp *http.Response, obj interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return &api.Error{Code: resp.StatusCode, Message: err.Error()}
	}
	logger.Debugf("Response body: %s", body)

	switch {
	case resp.StatusCode == http.StatusAccepted:
		logger.Debug("Request accepted")
	case resp.StatusCode == http.StatusNoContent:
		logger.Debug("Request suceeded. Response has no content")
	default:
		if resp.StatusCode > 399 && resp.StatusCode < 600 {
			var apiErr api.Error
			err = json.Unmarshal(body, &apiErr)
			if err != nil {
				// not json. 404s etc.
				return &api.Error{
					Code:    resp.StatusCode,
					Message: string(body),
				}
			}
			return &apiErr
		}
		err = json.Unmarshal(body, obj)
		if err != nil {
			return &api.Error{
				Code:    resp.StatusCode,
				Message: err.Error(),
			}
		}
	}
	return nil
}

func (c *defaultClient) handleStreamResponse(resp *http.Response, handler responseDecoder) error {
	if resp.StatusCode > 399 && resp.StatusCode < 600 {
		return c.handleResponse(resp, nil)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &api.Error{
			Code:    resp.StatusCode,
			Message: "expected streaming response with code 200",
		}
	}

	dec := json.NewDecoder(resp.Body)
	for {
		err := handler(dec)
		if err == io.EOF {
			// we need to check trailers
			break
		}
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	errTrailer := resp.Trailer.Get("X-Stream-Error")
	if errTrailer != "" {
		return &api.Error{
			Code:    500,
			Message: errTrailer,
		}
	}
	return nil
}