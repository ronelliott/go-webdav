package webdav

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/ronelliott/go-digest-auth-transport"
)

const (
	ClientOptionSSLNoVerify = "ssl-no-verify"
	DepthValueInfinity      = "infinity"
	HttpHeaderDepth         = "Depth"
	HttpHeaderDestination   = "Destination"
	HttpMethodCopy          = "COPY"
	HttpMethodMove          = "MOVE"
	HttpMethodMkcol         = "MKCOL"
	HttpMethodPropfind      = "PROPFIND"
	HttpUserAgentHeader     = "User-Agent"
	HttpUserAgentValue      = "GoWebDAV/0.1.0"
)

// The client implementation
type Client struct {
	// The base url for this Client
	BaseURL string

	// The http client for this client
	Client *http.Client

	// The password for this Client
	Password string

	// The root path for this Client
	RootPath string

	// The username for this Client
	Username string
}

// Create a new webdav client from the given resource
func NewClient(resource string) (*Client, error) {
	parsed, err := url.Parse(resource)

	if err != nil {
		return nil, err
	}

	path := parsed.Path
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	client := &Client{
		BaseURL:  fmt.Sprintf("%s://%s%s", parsed.Scheme, parsed.Host, path),
		RootPath: path,
	}

	var username string
	var password string

	if parsed.User != nil {
		username = parsed.User.Username()
		password, _ = parsed.User.Password()
	}

	sslNoVerifyRaw := parsed.Query().Get(ClientOptionSSLNoVerify)

	if sslNoVerifyRaw == "" {
		sslNoVerifyRaw = "false"
	}

	sslNoVerify, err := strconv.ParseBool(sslNoVerifyRaw)

	if err != nil {
		return nil, err
	}

	client.Client = dat.NewTransport(username, password, !sslNoVerify).Client()
	jar, err := cookiejar.New(nil)

	if err != nil {
		return nil, err
	}

	client.Client.Jar = jar

	// Try the connection, mainly to set any cookies required
	_, err = client.Root().Exists()

	if err != nil {
		return nil, err
	}

	return client, nil
}

// -----------------------------------------------------------------------------
// HTTP methods
// -----------------------------------------------------------------------------

// Run an copy request at the given resource
func (client *Client) copy(resource, destination string, depth int) (*http.Response, error) {
	req, err := client.make(HttpMethodCopy, resource, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set(HttpHeaderDepth, client.depth(depth))
	req.Header.Set(HttpHeaderDestination, client.url(destination))
	return client.run(req)
}

// Run an delete request at the given resource
func (client *Client) delete(resource string) (*http.Response, error) {
	return client.request(http.MethodDelete, resource, nil)
}

// Run an get request at the given resource
func (client *Client) get(resource string) (*http.Response, error) {
	return client.request(http.MethodGet, resource, nil)
}

// Run an head request at the given resource
func (client *Client) head(resource string) (*http.Response, error) {
	return client.request(http.MethodHead, resource, nil)
}

// Run an mkcol request at the given resource
func (client *Client) mkcol(resource string) (*http.Response, error) {
	return client.request(HttpMethodMkcol, resource, nil)
}

// Run an move request at the given resource
func (client *Client) move(resource, destination string) (*http.Response, error) {
	req, err := client.make(HttpMethodMove, resource, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set(HttpHeaderDepth, DepthValueInfinity)
	req.Header.Set(HttpHeaderDestination, client.url(destination))
	return client.run(req)
}

// Run an propfind request at the given resource
func (client *Client) propfind(resource string, depth int) (*http.Response, error) {
	req, err := client.make(HttpMethodPropfind, resource, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set(HttpHeaderDepth, client.depth(depth))
	return client.run(req)
}

// Run an put request at the given resource
func (client *Client) put(resource string, data []byte) (*http.Response, error) {
	return client.request(http.MethodPut, resource, data)
}

// -----------------------------------------------------------------------------
// Helper methods
// -----------------------------------------------------------------------------

// Returns the string form of the given depth value
func (client *Client) depth(depth int) string {
	depthValue := DepthValueInfinity

	if depth == 0 || depth == 1 {
		depthValue = fmt.Sprintf("%d", depth)
	}

	return depthValue
}

// Returns true if the response status is an error, or false otherwise
func (client *Client) isResponseError(res *http.Response) bool {
	return res.StatusCode >= http.StatusBadRequest
}

// Create and a request with the given method, resource and body
func (client *Client) make(method, resource string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, client.url(resource), bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	req.Header.Set(HttpUserAgentHeader, HttpUserAgentValue)
	return req, nil
}

// Get the location for the given response
func (client *Client) location(res *Response) string {
	return res.Location(client.RootPath)
}

// Read and return the data from the given response
func (client *Client) read(res *http.Response) ([]byte, error) {
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

// Create and run a request with the given method, resource and body, returning
// the data or any error encountered
func (client *Client) request(method, resource string, body []byte) (*http.Response, error) {
	req, err := client.make(method, resource, body)

	if err != nil {
		return nil, err
	}

	return client.run(req)
}

// Run the given request, returning the data or error
func (client *Client) run(req *http.Request) (*http.Response, error) {
	res, err := client.Client.Do(req)

	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		err = NotFoundError{req.URL.String()}
	case http.StatusUnauthorized:
		err = UnauthorizedError{req.URL.String()}

	default:
		if res.StatusCode >= http.StatusBadRequest {
			err = UnknownError{req.URL.String()}
		}
	}

	return res, err
}

// Get the full url for the given resource
func (client *Client) url(resource string) string {
	base := client.BaseURL

	if strings.HasPrefix(resource, client.RootPath) {
		resource = strings.Replace(resource, client.RootPath, "", 1)
	}

	if !strings.HasSuffix(base, "/") && !strings.HasPrefix(resource, "/") {
		base += "/"
	}

	return base + resource
}

// -----------------------------------------------------------------------------
// Public methods
// -----------------------------------------------------------------------------

// Get the collection resource at the given location
func (client *Client) Collection(location string) *Collection {
	return NewCollection(client, location)
}

// Get the entry resource at the given location
func (client *Client) Entry(location string) *Entry {
	return NewEntry(client, location)
}

// Get the root collection resource
func (client *Client) Root() *Collection {
	return client.Collection("/")
}
