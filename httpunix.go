// Package httpunix provides a HTTP transport (net/http.RoundTripper)
// that uses Unix domain sockets instead of HTTP.
//
// This is useful for non-browser connections within the same host, as
// it allows using the file system for credentials of both client
// and server, and guaranteeing unique names.
//
// The URLs look like this:
//
//     http+unix://unix:LOCATION:/PATH_ETC
//
// where LOCATION is the file system path to the unix socket file,
// and PATH_ETC follow normal http: scheme conventions.
package httpunix

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"
)

// Scheme is the URL scheme used for HTTP over UNIX domain sockets.
const Scheme = "http+unix"

// Transport is a http.RoundTripper that connects to Unix domain
// sockets.
type Transport struct {
	DialTimeout           time.Duration
	RequestTimeout        time.Duration
	ResponseHeaderTimeout time.Duration
}

// RoundTrip executes a single HTTP transaction. See
// net/http.RoundTripper.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL == nil {
		return nil, errors.New("http+unix: nil Request.URL")
	}
	if req.URL.Scheme != Scheme {
		return nil, errors.New("unsupported protocol scheme: " + req.URL.Scheme)
	}
	if req.URL.Host != "unix" {
		return nil, errors.New("http+unix: invalid Host in request URL")
	}

	parts := strings.Split(req.URL.Path, ":")
	if len(parts) != 2 {
		return nil, errors.New("http+unix: Invalid URL format")
	}

	unixSocketFile := parts[0]
	urlPath := parts[1]

	// change path during the request
	path := req.URL.Path
	defer func() {
		req.URL.Path = path
	}()
	req.URL.Path = urlPath

	c, err := net.DialTimeout("unix", unixSocketFile, t.DialTimeout)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(c)
	if t.RequestTimeout > 0 {
		c.SetWriteDeadline(time.Now().Add(t.RequestTimeout))
	}
	if err := req.Write(c); err != nil {
		return nil, err
	}
	if t.ResponseHeaderTimeout > 0 {
		c.SetReadDeadline(time.Now().Add(t.ResponseHeaderTimeout))
	}
	resp, err := http.ReadResponse(r, req)
	return resp, err
}
