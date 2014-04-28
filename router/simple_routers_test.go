package router

import (
	"bytes"
	"github.com/millken/falcore"
	"net/http"
	"regexp"
	"testing"
)

type SimpleFilter int

func (sf SimpleFilter) FilterRequest(req *falcore.Request) *http.Response {
	sf = -sf
	return nil
}

func validGetRequest() (req *falcore.Request) {
	tmp, _ := http.NewRequest("GET", "/hello", bytes.NewBuffer(make([]byte, 0)))
	req, _ = falcore.TestWithRequest(tmp, falcore.NewRequestFilter(func(req *falcore.Request) *http.Response { return nil }), nil)
	return
}

func TestRegexpRoute(t *testing.T) {
	r := new(RegexpRoute)

	var sf1 SimpleFilter = 1
	r.Filter = sf1
	r.Match = regexp.MustCompile(`one`)

	if r.MatchString("http://tester.com/one") != sf1 {
		t.Errorf("Failed to match regexp")
	}
	if r.MatchString("http://tester.com/two") != nil {
		t.Errorf("False regexp match")
	}

}

func TestHostRouter(t *testing.T) {
	hr := NewHostRouter()

	var sf1 SimpleFilter = 1
	var sf2 SimpleFilter = 2
	hr.AddMatch("www.ngmoco.com", sf1)
	hr.AddMatch("developer.ngmoco.com", sf2)

	req := validGetRequest()
	req.HttpRequest.Host = "developer.ngmoco.com"

	filt := hr.SelectPipeline(req)
	if filt != sf2 {
		t.Errorf("Host router didn't get the right pipeline")
	}

	req.HttpRequest.Host = "ngmoco.com"
	filt = hr.SelectPipeline(req)
	if filt != nil {
		t.Errorf("Host router got currently unsupported fuzzy match so you should update this test")
	}
}
