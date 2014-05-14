package filter

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"
)

// Holds all the information about how to lookup and
// connect to an upstream.
type UpstreamTransport struct {
	DNSCacheDuration time.Duration

	host string
	port int

	tcpaddr          *net.TCPAddr
	tcpaddrCacheTime time.Time

	transport *http.Transport
	timeout   time.Duration
}

// transport is optional.  We will override Dial
func NewUpstreamTransport(host string, port int, timeout time.Duration, transport *http.Transport) *UpstreamTransport {
	ut := &UpstreamTransport{
		host:      host,
		port:      port,
		timeout:   timeout,
		transport: transport,
	}
	ut.DNSCacheDuration = 15 * time.Minute

	if ut.transport == nil {
		ut.transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		ut.transport.MaxIdleConnsPerHost = 15
	}

	ut.transport.Dial = func(n, addr string) (c net.Conn, err error) {
		return ut.dial(n, addr)
	}

	return ut
}

func (t *UpstreamTransport) dial(n, a string) (c net.Conn, err error) {
	deadline := time.Now().Add(t.timeout)
	c, err = net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", t.host, t.port), t.timeout)
	if err != nil {
		return
	}
	c.SetDeadline(deadline)
	return
}
