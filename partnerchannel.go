package partnerchannel

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const (
	baseURL   = "https://my.partnerchannel.eu/api/v1/"
	mediaType = "application/json"
	userAgent = "PartnerChannel GO SDK"
)

// Client represents the client which interfaces with the API
type Client struct {
	httpClient   *http.Client
	apiToken     string
	baseURL      *url.URL
	sessionToken string
	Debug        bool

	// Endpoints
	Account Account
	Acme    Acme
}

// genericRequest represents the structure of the default request format
type genericRequest struct {
	Data     interface{}
	Metadata *MetadataRequest
}

// genericResponse represents the structure of the default response format
type genericResponse struct {
	Metadata *MetadataResponse
}

// NewClient creates a new cliente with the given apiToken and httpClient (Uses http.DefaultClient by default)
func NewClient(apiToken string, httpClient *http.Client) (*Client, error) {
	client := &Client{}

	if httpClient == nil {
		client.httpClient = http.DefaultClient
	}

	var err error

	client.apiToken = apiToken
	client.baseURL, err = url.Parse(baseURL)

	// Setup the endpoints
	client.Account = &accountEndpoint{client: client}
	client.Acme = &acmeEndpoint{client: client}

	return client, err
}

// ChangeBaseURL changes the Base-URL
func (c *Client) ChangeBaseURL(newBaseUrl *url.URL) {
	c.baseURL = newBaseUrl
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client. Relative URLs should always be specified without a preceding slash. If specified, the
// value pointed to by body is JSON encoded and included in as the request body.
func (c *Client) NewRequest(method, urlStr string, query url.Values, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	if query != nil {
		u.RawQuery = query.Encode()
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("X-API-Key", c.apiToken)

	if c.sessionToken != "" {
		req.Header.Add("X-Session-ID", c.sessionToken)
	}

	if c.Debug {
		dump, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("%s", dump)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is JSON decoded and stored in the value
// pointed to by v, or returned as an error if an API error has occurred. If v implements the io.Writer interface,
// the raw response will be written to v, without attempting to decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Debug {
		dump, err := httputil.DumpResponse(resp, true)
		if err != nil {
			log.Println(err)
		}

		fmt.Printf("%s", dump)
	}

	// The PartnerChannel returns 200 in every case.
	// If we do not recieve an 200 we need to fail instantly
	if resp.StatusCode != http.StatusOK {
		return resp, errors.New(resp.Status)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	if v != nil {
		err := json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return nil, err
		}
	}

	return resp, err
}
