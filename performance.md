---
layout: default
title: Measuring Falcore Performance
---

The pipeline tracks performance timing on every stage in the pipeline and exposes that information in a final callback stage after the request/response has been completed.  This allows fine grained tracking of performance of requests as they traverse the pipeline including the IO read and write timings.  This information may be ignored or captured via the callback to do with as you please.

Here’s a simple working example of a Falcore server that illustrates some of the features.

```go
package main
import (
        "falcore"
        "http"
        "fmt"
        "time"
        "rand"
)
func main() {
        // create pipeline
        pipeline := falcore.NewPipeline()
        // add upstream pipeline stages
        var filter1 delayFilter
        pipeline.Upstream.PushBack(filter1)
        var filter2 helloFilter
        pipeline.Upstream.PushBack(filter2)
        // add request done callback stage
        pipeline.RequestDoneCallback = reqCB

        // create server on port 8000
        server := falcore.NewServer(8000, pipeline)
		
		server.CompletionCallback = func(req *falcore.Request, res *http.Response ) {
			req.Trace(res)
		}
		
        // start the server
        // this is normally blocking forever unless you send lifecycle commands
        if err := server.ListenAndServe(); err != nil {
                fmt.Println("Could not start server:", err)
        }
}
// Example filter to show timing features
type delayFilter int
func (f delayFilter) FilterRequest(req *falcore.Request) *http.Response {
        status := rand.Intn(2) // random status 0 or 1
        if status == 0 {
                time.Sleep(rand.Int63n(100e6)) // random sleep between 0 and 100 ms
        }
        req.CurrentStage.Status = byte(status) // set the status to produce a unique signature
        return nil
}
// Example filter that returns a Response
type helloFilter int
func (f helloFilter) FilterRequest(req *falcore.Request) *http.Response {
        return falcore.StringResponse(req.HttpRequest, 200, nil, "hello world!\n")
}
```

First, the pipeline is created and two trivial filter stages are added.  Then we add the request done callback. Finally, we create and start the server.  Using any HTTP client, we can make a few requests and see the output:

	2012/01/27 13:25:55 [TRAC] 81439859c8 [GET] localhost:8000/ S=200 Sig=2E23AE3C Tot=0.0003
	2012/01/27 13:25:55 [TRAC] 81439859c8 server.Init                    S=0 Tot=0.0002 %=65.70
	2012/01/27 13:25:55 [TRAC] 81439859c8 main.delayFilter               S=1 Tot=0.0000 %=0.97
	2012/01/27 13:25:55 [TRAC] 81439859c8 main.helloFilter               S=0 Tot=0.0000 %=3.56
	2012/01/27 13:25:55 [TRAC] 81439859c8 server.ResponseWrite           S=0 Tot=0.0001 %=24.60
	2012/01/27 13:25:55 [TRAC] 81439859c8 Overhead                       S=0 Tot=0.0000 %=5.18
	2012/01/27 13:25:56 [TRAC] 8181861943 [GET] localhost:8000/ S=Sig=2E23AE3C Tot=0.0002
	2012/01/27 13:25:56 [TRAC] 8181861943 server.Init                    S=0 Tot=0.0001 %=49.11
	2012/01/27 13:25:56 [TRAC] 8181861943 main.delayFilter               S=1 Tot=0.0000 %=1.18
	2012/01/27 13:25:56 [TRAC] 8181861943 main.helloFilter               S=0 Tot=0.0000 %=1.78
	2012/01/27 13:25:56 [TRAC] 8181861943 server.ResponseWrite           S=0 Tot=0.0001 %=38.46
	2012/01/27 13:25:56 [TRAC] 8181861943 Overhead                       S=0 Tot=0.0000 %=9.47
	2012/01/27 13:25:56 [TRAC] 81a7d3c473 [GET] localhost:8000/ Sig=60AAA595 Tot=0.0944
	2012/01/27 13:25:56 [TRAC] 81a7d3c473 server.Init                    S=0 Tot=0.0001 %=0.11
	2012/01/27 13:25:56 [TRAC] 81a7d3c473 main.delayFilter               S=0 Tot=0.0941 %=99.71
	2012/01/27 13:25:56 [TRAC] 81a7d3c473 main.helloFilter               S=0 Tot=0.0000 %=0.01
	2012/01/27 13:25:56 [TRAC] 81a7d3c473 server.ResponseWrite           S=0 Tot=0.0001 %=0.14
	2012/01/27 13:25:56 [TRAC] 81a7d3c473 Overhead                       S=0 Tot=0.0000 %=0.03

There are three requests logged and the log output comes from the `request.Trace()` call.  You can clearly see our stages that we created and how long each took.

### Request Ids

Debugging production issues can be very difficult when the application is under heavy load.  Often the information necessary to resolve an issue is present in the log but lost in the noise of other requests.  Falcore tracks requests in the pipeline in two different ways to facilitate better debugging.  The first is that each request is given a unique ID when created.  Printing this value in your log messages allows you to easily grep the messages for a given request.  The output above shows the request ID to the right of the level “`[TRAC]`.”  

### Signatures

The second helpful feature is the request signature.  The signature is a unique identifier for each possible path through the pipeline.  Each pipeline stage has a status that defaults to zero, but if your stage changes the status, then Falcore will produce a different unique signature.  In the example above, the delayFilter randomly chooses to sleep and the choice is reflected in the status.  When you split your performance metrics based on the unique signatures, you can easily see the difference between the two request types (in this case, 0.2-0.3ms vs. 94.4ms).

Setting unique status codes for errors allows easy tracking of the occurrence of those errors based on their signatures.  We’ve found that this type of tracking helps isolate specific issues to certain types of requests allowing us to find a reproduce issues much more quickly than we would be able to do otherwise.
