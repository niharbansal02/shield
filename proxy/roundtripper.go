package proxy

import (
	"golang.org/x/net/http2"
	"net"
	"net/http"
	"time"

	"github.com/odpf/salt/log"
)

type h2cTransportWrapper struct {
	transport *http.Transport
	log       log.Logger
}

func (t *h2cTransportWrapper) RoundTrip(req *http.Request) (*http.Response, error) {
	// we need to apply errors if it failed in Director
	if err, ok := req.Context().Value(CtxRequestErrorKey).(error); ok {
		return nil, err
	}
	t.log.Debug("proxy request", "host", req.URL.Host, "path", req.URL.Path,
		"scheme", req.URL.Scheme, "protocol", req.Proto)
	return t.transport.RoundTrip(req)
}

func NewH2cRoundTripper(log log.Logger) http.RoundTripper {

	transport := &http.Transport{
		DialContext: (&net.Dialer{ // use DialContext here
			Timeout:   10 * time.Second,
			KeepAlive: 1 * time.Minute,
		}).DialContext,
	}

	if err := http2.ConfigureTransport(transport); err != nil {
		log.Error(err.Error())
	}

	return &h2cTransportWrapper{
		transport: transport,
		log:       log,
	}
}
