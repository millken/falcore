package responder

import (
	"github.com/millken/falcore"
	"net/http"
)

// A 302 redirect response
func RedirectResponse(req *http.Request, url string) *http.Response {
	h := make(http.Header)
	h.Set("Location", url)
	return falcore.SimpleResponse(req, 302, h, 0, nil)
}
