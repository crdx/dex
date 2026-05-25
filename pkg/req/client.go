// Package req provides a lightweight chainable JSON HTTP client built on net/http.
package req

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

// Client is an HTTP client with optional debug logging and default headers.
type Client struct {
	client  *http.Client
	debug   bool
	headers map[string]string
}

// NewClient creates a new Client.
func NewClient() *Client {
	return &Client{
		client:  http.DefaultClient,
		headers: make(map[string]string),
	}
}

// SetDebug enables debug logging of requests and responses to stderr.
func (self *Client) SetDebug(debug bool) *Client {
	self.debug = debug
	return self
}

// SetBearerAuth sets the Authorization header to a bearer token for all requests.
func (self *Client) SetBearerAuth(token string) *Client {
	self.headers["Authorization"] = "Bearer " + token
	return self
}

// R creates a new Request from this client.
func (self *Client) R() *Request {
	return &Request{
		client: self,
	}
}

// Request is a chainable request builder.
type Request struct {
	client        *Client
	body          any
	queryParams   map[string]string
	successResult any
	errorResult   any
}

// SetBody sets the JSON request body.
func (self *Request) SetBody(body any) *Request {
	self.body = body
	return self
}

// SetQueryParams sets URL query parameters.
func (self *Request) SetQueryParams(params map[string]string) *Request {
	self.queryParams = params
	return self
}

// SetSuccessResult sets the target for JSON unmarshalling on 2xx responses.
func (self *Request) SetSuccessResult(result any) *Request {
	self.successResult = result
	return self
}

// SetErrorResult sets the target for JSON unmarshalling on non-2xx responses.
func (self *Request) SetErrorResult(result any) *Request {
	self.errorResult = result
	return self
}

// Send executes the request and returns a Response.
func (self *Request) Send(method string, endpoint string) (*Response, error) {
	if self.queryParams != nil {
		values := url.Values{}
		for key, value := range self.queryParams {
			values.Set(key, value)
		}
		endpoint += "?" + values.Encode()
	}

	var body io.Reader
	if self.body != nil {
		jsonBody, err := json.Marshal(self.body)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(jsonBody)
	}

	httpRequest, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	for key, value := range self.client.headers {
		httpRequest.Header.Set(key, value)
	}

	if self.body != nil {
		httpRequest.Header.Set("Content-Type", "application/json")
	}

	if self.client.debug {
		fmt.Fprintf(os.Stderr, "%s %s\n", method, endpoint)
	}

	httpResponse, err := self.client.client.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer func() { _ = httpResponse.Body.Close() }()

	responseBody, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	if self.client.debug {
		fmt.Fprintf(os.Stderr, "%d %s\n", httpResponse.StatusCode, string(responseBody))
	}

	response := &Response{
		StatusCode: httpResponse.StatusCode,
		body:       responseBody,
	}

	if response.IsSuccessState() && self.successResult != nil {
		if err := json.Unmarshal(responseBody, self.successResult); err != nil {
			return response, err
		}
	}

	if !response.IsSuccessState() && self.errorResult != nil {
		// Best-effort unmarshal of error body. If it fails, the caller can still check the status.
		_ = json.Unmarshal(responseBody, self.errorResult)
	}

	return response, nil
}

// Response wraps an HTTP response.
type Response struct {
	StatusCode int
	body       []byte
}

// IsSuccessState returns true if the status code is 2xx.
func (self *Response) IsSuccessState() bool {
	return self.StatusCode >= 200 && self.StatusCode < 300
}
