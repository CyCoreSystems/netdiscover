package discover

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
)

// ParseIPv4FromBody reads from a reader (such as an http.Body) and tries to
// parse the IPv4 IP address contained therein.
func ParseIPv4FromBody(in io.Reader) (net.IP, error) {

	data, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	if len(data) < 8 || len(data) > 42 {
		return nil, fmt.Errorf("invalid response: %s", string(data))
	}

	ip := net.ParseIP(string(data))
	if ip == nil {
		return nil, fmt.Errorf("failed to parse IP: %s", string(data))
	}

	return ip, nil
}

// StandardIPFromHTTP queries an HTTP URL and returns the IP address contained within its response.  Headers are optional.
func StandardIPFromHTTP(url string, headers map[string]string) (net.IP, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to construct HTTP request: %v", err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("non-2XX response: (%d) %s", resp.StatusCode, resp.Status)
	}

	return ParseIPv4FromBody(resp.Body)
}

// StandardHostnameFromHTTP queries an HTTP URL and returns the hostname contained within its response.  Headers are optional.
func StandardHostnameFromHTTP(url string, headers map[string]string) (string, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to construct HTTP request: %v", err)
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return "", fmt.Errorf("non-2XX response: (%d) %s", resp.StatusCode, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}
	if len(body) < 4 {
		return "", fmt.Errorf("hostname implausibly short: %s", string(body))
	}

	return string(body), nil
}
