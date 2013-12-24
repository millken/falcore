---
layout: default
title: Falcore
---
## What is Falcore?

Falcore is a framework for constructing high performance, modular HTTP servers in [Go][go] (sometimes called Golang).

Falcore's architecture consists of an upstream pipeline where stages are run in sequence until a stage returns an HTTP Response.  There is a separate downstream pipeline which runs against the response before it is returned to the client.  Any stage may be arbitrarily complex, but you'll get extra benefits from breaking your application up into separate modules spread out over stages.

An example request might encounter several filters while traversing a pipeline:

* A authorization filter that performs application specific auth, then either responds with an error or appends auth info to the request's context
* A router that looks at the request and chooses one of several filters to handle
* A filter that evaluates cache headers to see the if the response can be served from a local cache
* A filter that peforms the core functionality of the application, such as rendering a JSON dictionary for an API call
* A downstream filter that evaluates cache headers and response to determine if the response should be added into the local cache
* A downstream filter that implements etag/if-none-match comparison
* A downstream filter that implements gzip/deflate compression if appropriate and possible

One central benefit to pipelining approach is that it allows for code re-use.  Included in the repositories are filters for performing common HTTP server tasks such as response compression, etag matching, and serving files from disk.  Each of these filters, and any filters you write yourself, can easily be dropped into any Falcore pipeline.  Simply compose the features you want and start your server.

You can [read the full documentation on Godoc.org](http://godoc.org/github.com/fitstar/falcore).

## Features

* Modular and flexible design
* [Hot restart hooks](hot_restart.html) for zero-downtime deploys
* [Builtin performance statistics](performance.html) framework
* Builtin logging framework
* [Compatible](compatibility.html) with `net/http` and Google App Engine

## Using Falcore

Falcore is a filter pipeline based HTTP server library.  You can build arbitrarily complicated HTTP services by chaining just a few simple components:
	
`RequestFilters` are the core component.  A request filter takes a request and returns a response or nil.  Request filters can modify the request as it passes through.

`ResponseFilters` can modify a response on its way out the door.  An example response filter, `compression_filter`, is included.  It applies `deflate` or `gzip` compression to the response if the request supplies the proper headers.

`Pipelines` form one of the two logic components.  A pipeline contains a list of `RequestFilters` and a list of `ResponseFilters`.  A request is processed through the request filters, in order, until one returns a response.  It then passes the response through each of the response filters, in order.  A pipeline is a valid `RequestFilter`.

`Routers` allow you to conditionally follow different pipelines.  A router chooses from a set of pipelines.  A few basic routers are included, including routing by hostname or requested path.  You can implement your own router by implementing `falcore.Router`.  `Routers` are not `RequestFilters`, but they can be put into pipelines.

See the `examples` directory for usage examples.

## Getting falcore

Install with `go get github.com/fitstar/falcore`.

## HTTPS

To use falcore to serve HTTPS, simply call `ListenAndServeTLS` instead of `ListenAndServe`.  If you want to host SSL and nonSSL out of the same process, simply create two instances of `falcore.Server`.  You can give them the same pipeline or share pipeline components.

## Building

Falcore is currently targeted at Go 1.0.  If you're still using Go r.60.x, you can get the last working version of falcore for r.60 using the tag `last_r60`.

## Maintainers

* [Dave Grijalva](http://www.github.com/dgrijalva)
* [Scott White](http://www.github.com/smw1218)

## Contributors

* [Graham Anderson](http://www.github.com/gnanderson)
* [Amir Mohammad Saied](http://github.com/amir)
* [James Wynn](https://github.com/jameswynn)
* [Jonathan Rudenberg](https://github.com/titanous)



[go]: http://www.golang.org
