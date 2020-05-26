package oktalib

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type doRequestInput struct {
	Dest      string // The destination address. It has to be hostname
	Dport     string // The destination address port
	TimeOut   int
	Method    string // The req method, POST/PATCH etc...
	Route     string // The route, by default its "/" it can be "/api"
	Headers   map[string]string
	Payload   interface{}
	CookieJar http.CookieJar
}

// The chain of request need a custom transport layer to work.
func doRequest(i *doRequestInput) (*http.Response, error) {
	// Initialization
	var (
		resp          *http.Response
		cert          tls.Certificate
		certlist      []tls.Certificate
		err           error
		encodepayload []byte
		addr          string
		portinfo      string
		ebody         *bytes.Reader
	)
	client := &http.Client{
		Jar: i.CookieJar,
	}

	certlist = append(certlist, cert)

	// Load our CA certificate

	tlsConfig := &tls.Config{
		Certificates: certlist,
	}
	client.Transport = &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: tlsConfig,
	}

	if i.TimeOut == 0 {
		client.Timeout = time.Duration(20000) * time.Millisecond
	} else {
		client.Timeout = time.Duration(i.TimeOut) * time.Millisecond
	}
	encodepayload, err = json.Marshal(i.Payload)
	if len(i.Route) > 0 {
		if string(i.Route[0]) != "/" {
			i.Route = "/" + i.Route
		}
	}
	if i.Dport == "" {
		portinfo = ""
	} else {
		portinfo = ":" + i.Dport
	}
	addr = i.Dest + portinfo + i.Route
	ebody = bytes.NewReader(encodepayload)
	req, err := http.NewRequest(i.Method, addr, ebody)
	if err != nil {
		return resp, fmt.Errorf("Error making new request %v", err)
	}
	for k, v := range i.Headers {
		req.Header.Set(k, v)
	}
	resp, err = client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("client do error %v", err)
	}
	return resp, nil
}
