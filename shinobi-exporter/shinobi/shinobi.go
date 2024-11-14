package shinobi

import (
//        "io"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	DefaultClientTimeout  = 5 * time.Second
	DefaultRequestTimeout = 2500 * time.Millisecond
//        DefaultRequestTimeout = 5
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Server struct {
	endpoint string

	client HTTPClient
}

type Group string
type Groups []Group

func NewServerDefault(endpoint, key string, insecure bool) (*Server, error) {
	uri, err := url.Parse(endpoint + "/" + key)
	if err != nil {
		return nil, fmt.Errorf("incorrect server endpoint addr: %w", err)
	}

	client := &http.Client{
		Timeout: DefaultClientTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}

	return &Server{
		endpoint: uri.String(),
		client:   client,
	}, nil
}

func (s *Server) sendGetRequest(ctx context.Context, resource string, result any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.endpoint+resource, nil)
	if err != nil {
		return fmt.Errorf("build resource %s request error: %w", resource, err)
	}

	res, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("do resource %s request  error: %w", resource, err)
	}

	if res.StatusCode != http.StatusOK {
		return ErrorUnexpecterAPIResponseStatus(res.StatusCode)
	}

	defer res.Body.Close()

//                bodyBytes, err := io.ReadAll(res.Body)
//                if err != nil {
//                   fmt.Errorf("%w", err)
//                }
//                bodyString := string(bodyBytes)
 //               fmt.Printf("%s", bodyString)

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("response resource %s decode error: %w", resource, err)
	}

	return nil
}
