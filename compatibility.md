---
layout: default
title: Falcore Compatibility
---

### What's up with compatibility?

While Go makes it simple to create a web server using only their excellent HTTP libraries, for large and long-lived projects, it’s important to have the sort of features that Falcore provides.  The `falcore.Server` provides several critical hooks into the application lifecycle for hot restart and other production application lifecycle management tasks.  Composible pipelines make it far easier to share functionality between applications, both internally or publicly.

All that said, we want falcore to integrate well with the standard library's `net/http` package.  Falcore is fully compatible with `http.Server` and `http.Handler`.  This means falcore is totally compatible with [Google App Engine](https://developers.google.com/appengine/docs/go/overview).

### http.Response vs http.ResponseWriter

The standard library's server functionality makes use of the `http.ResponseWriter` interface.  While this interface is totally fine for just writing a response to a socket, it's not conducive to the multi-stage pre-and-post-processing you might want to do with falcore.  Luckily, `http.Response`, also part of the standard library, is a struct describing all the relevant details of an HTTP response.

Falcore makes use of `http.Response` instead of `http.ResponseWriters` to make downstream processing simpler.  Falcore also provides several helper methods for generating these response objects, so you don't have to figure out detaily bits yourself (See:`falcore.StringResponse`, `falcore.JSONResponse`, and friends).

### Using an http.Handler with falcore

The `falcore/filter` package includes `filter.HandlerFilter`, which will use your `http.Handler` as a `falcore.ResponseFilter`.  The one caveat here is that, since an `http.Handler` is expected to always return a response, we can't implement include the fallthrough behavior of standard response filters.  This means any response filters after a `HandlerFilter` will be unreachable.

You could, however, use a `falcore.Router` to extract your routing behavior and simplify your `Handler`s.

### Using falcore with http.Server

`falcore.Server` is a valid `http.Handler`.  To use this functionality, simply create your `falcore.Server` instance, configure it as you normally would, and pass it off to your `http.Server` as a handler instead of calling `ListenAndServe`.  Falcore will still track stats for all stages in the pipeline and call `CompletionCallback` after each request.

Unfortunately, features like hot restart can't be provided if falcore doesn't control the socket listener.

### Using both

You can totally use both of these techniques.  For example, you may want to add some post processing pipeline stages to a Google App Engine application.  Simply follow both of these techniques and insert anything else you like into the falcore pipeline.  Happy coding!