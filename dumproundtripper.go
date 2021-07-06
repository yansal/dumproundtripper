package dumproundtripper

import (
	"log"
	"net/http"
	"net/http/httputil"
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
	body := true
	reqbody, err := httputil.DumpRequestOut(req, body)
	if err != nil {
		return nil, err
	}
	rt.logger.Printf("%s", reqbody)
	resp, err := rt.wrapped.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	respbody, err := httputil.DumpResponse(resp, body)
	if err != nil {
		return nil, err
	}
	rt.logger.Printf("%s", respbody)
	return resp, nil
}
