package responder

import (
	"bytes"
	"encoding/json"
	"github.com/fitstar/falcore"
	"net/http"
)

// Generate an http.Response by json encoding body using
// the standard library's json.Encoder.  error will be nil
// unless json encoding fails.
func JSONResponse(req *http.Request, status int, headers http.Header, body interface{}) (*http.Response, error) {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(body); err != nil {
		return nil, err
	}

	if headers == nil {
		headers = make(http.Header)
	}
	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/json")
	}

	return falcore.SimpleResponse(req, status, headers, int64(buf.Len()), buf), nil
}

// Streaming version of JSONResponse.  JSON encoding is unbuffered and transfer-encoding will be chunked.
// Errors encountered during encoding will be logged, but are not returned
func StreamingJSONResponse(req *http.Request, status int, headers http.Header, body interface{}) *http.Response {
	if headers == nil {
		headers = make(http.Header)
	}
	if headers.Get("Content-Type") == "" {
		headers.Set("Content-Type", "application/json")
	}

	pW, res := falcore.PipeResponse(req, status, headers)
	go func() {
		if err := json.NewEncoder(pW).Encode(body); err != nil {
			Error("Error encoding JSON: %v", err)
		}
		pW.Close()
	}()

	return res
}
