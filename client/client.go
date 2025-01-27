package client

import (
	"net"
	"net/http"
	"os"
	"time"
)

var client *http.Client
var transport *http.Transport

const envCrawlerUA = "CRAWLER_USER_AGENT"

type myTransport struct {
	Transport http.RoundTripper
}

func (t *myTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if os.Getenv(envCrawlerUA) != "" {
		req.Header.Set("User-Agent", envCrawlerUA)
	}
	return t.Transport.RoundTrip(req)
}

func init() {
	transport = http.DefaultTransport.(*http.Transport).Clone()
	transport.DialContext = (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		Control:   control,
	}).DialContext

	client = &http.Client{
		Transport: &myTransport{Transport: transport},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

}

func Get(url string) (*http.Response, error) {
	return (&http.Client{
		Transport: &myTransport{Transport: transport},
	}).Get(url)
}
