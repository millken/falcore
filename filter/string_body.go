package filter

import (
	"bytes"
	"github.com/millken/falcore"
	"github.com/millken/falcore/utils"
	"io"
	"net/http"
	"strings"
)

// Keeps the body of a request in a string so it can be re-read at each stage of the pipeline
// implements io.ReadCloser to match http.Request.Body
type StringBody struct {
	BodyBuffer *bytes.Reader
	bpe        *utils.BufferPoolEntry
}

// Injests the body of the request into a buffer so it can be read more than
// once.
type StringBodyFilter struct {
	pool *utils.BufferPool
}

func NewStringBodyFilter() *StringBodyFilter {
	sbf := &StringBodyFilter{}
	sbf.pool = utils.NewBufferPool(100, 1024)
	return sbf
}
func (sbf *StringBodyFilter) FilterRequest(request *falcore.Request) *http.Response {
	req := request.HttpRequest
	// This caches the request body so that multiple filters can iterate it
	if req.Method == "POST" || req.Method == "PUT" {
		sb, err := sbf.readRequestBody(req)
		if sb == nil || err != nil {
			request.CurrentStage.Status = 3 // Skip
			falcore.Debug("%s No Req Body or Ignored: %v", request.ID, err)
		}
	} else {
		request.CurrentStage.Status = 1 // Skip
	}
	return nil
}

// reads the request body and replaces the buffer with self
// returns nil if the body is multipart and not replaced
func (sbf *StringBodyFilter) readRequestBody(r *http.Request) (sb *StringBody, err error) {
	ct := r.Header.Get("Content-Type")
	// leave it on the buffer if we're multipart
	if strings.SplitN(ct, ";", 2)[0] != "multipart/form-data" && r.ContentLength > 0 {
		sb = &StringBody{}
		const maxFormSize = int64(10 << 20) // 10 MB is a lot of text.
		sb.bpe = sbf.pool.Take(io.LimitReader(r.Body, maxFormSize+1))

		// There shouldn't be a null byte so we should get EOF
		b, e := sb.bpe.Br.ReadBytes(0)
		if e != nil && e != io.EOF {
			return nil, e
		}
		sb.BodyBuffer = bytes.NewReader(b)
		r.Body.Close()
		r.Body = sb
		return sb, nil
	}
	return nil, nil // ignore
}

// Returns a buffer used in the FilterRequest stage to a buffer pool
// this speeds up this filter significantly by reusing buffers
func (sbf *StringBodyFilter) ReturnBuffer(request *falcore.Request) {
	if sb, ok := request.HttpRequest.Body.(*StringBody); ok {
		sbf.pool.Give(sb.bpe)
	}
}

// Insert this in the response pipeline to return the buffer pool for the request body
// If there is an appropriate place in your flow, you can call ReturnBuffer explicitly
func (sbf *StringBodyFilter) FilterResponse(request *falcore.Request, res *http.Response) {
	sbf.ReturnBuffer(request)
}

func (sb *StringBody) Read(b []byte) (n int, err error) {
	return sb.BodyBuffer.Read(b)
}

func (sb *StringBody) Close() error {
	// start over
	sb.BodyBuffer.Seek(0, 0)
	return nil
}
