package dumproundtripper

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

// New returns a new http.RoundTripper that dumps all request and response to
// logger.
func New(wrapped http.RoundTripper, logger *log.Logger) http.RoundTripper {
	return &dumproundtripper{
		logger:  logger,
		wrapped: wrapped,
	}
}

type dumproundtripper struct {
	logger  *log.Logger
	wrapped http.RoundTripper
}

func (rt *dumproundtripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		dumpbody       bool
		reqcontenttype = req.Header.Get("Content-Type")
	)
	if strings.HasPrefix(reqcontenttype, "application/json") ||
		strings.HasPrefix(reqcontenttype, "text/plain") {
		dumpbody = true
	}
	reqbody, err := httputil.DumpRequestOut(req, dumpbody)
	if err != nil {
		return nil, err
	}
	rt.logger.Printf("%s", reqbody)
	resp, err := rt.wrapped.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	dumpbody = false
	respcontenttype := resp.Header.Get("Content-Type")
	if strings.HasPrefix(respcontenttype, "application/json") ||
		strings.HasPrefix(respcontenttype, "text/plain") {
		dumpbody = true
	}

	respbody, err := httputil.DumpResponse(resp, dumpbody)
	if err != nil {
		return nil, err
	}
	rt.logger.Printf("%s", respbody)
	return resp, nil
}
