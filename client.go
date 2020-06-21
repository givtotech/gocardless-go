package gocardless

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

const (
	apiVersion     = "2015-07-06"
	baseLiveURL    = `https://api.gocardless.com/`
	baseSandboxURL = `https://api-sandbox.gocardless.com/`
)

// Client for interacting with the GoCardless Pro API
type Client struct {
	// AccessToken is the bearer token used to authenticate requests to the GoCardless API
	AccessToken string
	// RemoteURL is the address of the GoCardless API
	RemoteURL string
	// httpClient used for APi requests
	httpClient *http.Client
}

// NewClient instantiate a client struct with your access token and environment, then
// use the resource methods to access the API
func NewClient(accessToken string, env Environment) *Client {
	return NewClientWithHTTPClient(&http.Client{}, accessToken, env)
}

// NewClientWithHTTPClient instantiate a client struct with your access token and environment, then
// use the resource methods to access the API. Uses existing http.Client to allow customisations
func NewClientWithHTTPClient(hc *http.Client, accessToken string, env Environment) *Client {
	c := &Client{
		AccessToken: accessToken,
		httpClient:  hc,
	}

	switch env {
	case SandboxEnvironment:
		c.RemoteURL = baseSandboxURL
	case LiveEnvironment:
		c.RemoteURL = baseLiveURL
	default:
		log.Fatalf("Invalid environment %s, use one of (%s, %s)", env, SandboxEnvironment, LiveEnvironment)
	}
	return c
}

func (c *Client) makeRequest(ctx context.Context, path, method string, body, dst interface{}) error {
	req, err := c.newRequest(ctx, path, method, body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return errors.New("StatusTooManyRequests")
	}

	res := newResponse(resp)
	// bind response to struct
	return res.bind(dst)
}

func (c *Client) newRequest(ctx context.Context, path string, method string, body interface{}) (*http.Request, error) {
	if strings.ToUpper(method) == http.MethodPatch {
		return nil, errors.New(InvalidMethodError)
	}

	url := fmt.Sprintf("%s%s", c.RemoteURL, path)

	var bs []byte
	if body != nil {
		bs, _ = json.Marshal(body)
	}

	data := bytes.NewBuffer(bs)
	req, err := http.NewRequestWithContext(ctx, method, url, data)
	if err != nil {
		return nil, err
	}

	// set default headers
	c.setDefaultHeaders(req)

	if method == http.MethodPost {
		// Add Idempotency header key when creating a resouce
		// https://developer.gocardless.com/api-reference/#making-requests-idempotency-keys
		u, _ := uuid.NewV4()
		req.Header.Add("Idempotency-Key", u.String())
	}

	return req, nil
}

func (c *Client) setDefaultHeaders(req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	req.Header.Add("GoCardless-Version", apiVersion)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
}

func (c *Client) get(ctx context.Context, path string, dst interface{}) error {
	return c.makeRequest(ctx, path, http.MethodGet, nil, dst)
}

func (c *Client) post(ctx context.Context, path string, body, dst interface{}) error {
	return c.makeRequest(ctx, path, http.MethodPost, body, dst)
}

func (c *Client) put(ctx context.Context, path string, body, dst interface{}) error {
	return c.makeRequest(ctx, path, http.MethodPut, body, dst)
}

func (c *Client) delete(ctx context.Context, path string) error {
	return c.makeRequest(ctx, path, http.MethodDelete, nil, nil)
}
