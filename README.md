[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](http://godoc.org/github.com/agorman/httpunix)

Package httpunix provides a HTTP transport (net/http.RoundTripper) that uses Unix domain sockets instead of HTTP.

forked from https://github.com/tv42/httpunix

## Why

I forked this project because I wanted to provide a URL representation for http+unix addresses rather than have to register each unix socket programmatically.

## URLs

The URLs look like this:

```
http+unix://unix:LOCATION:/PATH_ETC
```

where LOCATION is the file system path to the unix socket file,
and PATH_ETC follow normal http: scheme conventions.

## Example

```
// This example shows handling all net/http requests for the
// http+unix URL scheme.
u := &httpunix.Transport{
	DialTimeout:           100 * time.Millisecond,
	RequestTimeout:        1 * time.Second,
	ResponseHeaderTimeout: 1 * time.Second,
}

// If you want to use http: with the same client:
t := &http.Transport{}
t.RegisterProtocol(httpunix.Scheme, u)
var client = http.Client{
	Transport: t,
}

resp, err := client.Get("http+unix://unix:/path/to/socket:/urlpath/as/seen/by/server")
if err != nil {
	log.Fatal(err)
}
buf, err := httputil.DumpResponse(resp, true)
if err != nil {
	log.Fatal(err)
}
fmt.Printf("%s", buf)
resp.Body.Close()
```
