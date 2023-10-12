package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

// parseProxyURI, configure the connection to the proxy, returns *http.Transport
// if parsing the URL data was correct.
func parseProxyURI() (*http.Transport, error) {
	proxySecret := os.Getenv("PROXY_URL")
	if proxySecret == "" {
		return &http.Transport{}, fmt.Errorf("PROXY_URL string is empty in .env file.")
	}

	proxyURL, err := url.Parse(proxySecret)
	if err != nil {
		return &http.Transport{}, fmt.Errorf("error while parsing URL of proxy: %w", err)
	}

	return &http.Transport{Proxy: http.ProxyURL(proxyURL)}, nil

}
