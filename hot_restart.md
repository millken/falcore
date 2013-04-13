---
layout: default
title: Falcore Hot Restart
---

## What is 'hot restart'?

There are many reasons you may need to restart a server process: config file changes, software changes, etc.  It sure would be swell to be able to do this without dropping requests mid-flight.  There are ways to accomplish this with things like nginx, but one of the design goals for falcore was to make sure it could be run right on the socket without needing things like nginx, or haproxy in front.  

In order to satisfy these two goals, we've provided what we call 'hot restart' functionality.  Most of the goods are simply example code, but falcore provides the necessary lifecycle management features to make it work.  To see this in action, check out `examples/hot_restart`.  You can actually copy/paste the important bits into your falcore app and it will work.  

## What happens when my app hot restarts?

In short, your process is replaced by a new process without dropping connections or closing the socket listener.  

For more detail, lets look at `examples/hot_restart`.  When you start the application, it will:

* setup some signal handlers so you can trigger hot restart from outside
* create a falcore server instance, which will
* open a socket listener on the specified port.  this is what waits for new connections

Once the app is running, if you send it `SIGHUP`, it will:

* use `syscall.ForkExec` to start a brand new instance of the process with
	* the falcore listener file descriptor kept open
	* an additional command line flag specifying the listener's file descriptor
* the child process will start up as usual, using the specified file pointer as the listener's file descriptor

At this point, we have a parent process (let's call it `A`) and a child process (`B`).  They are sharing the socket listener.  Once `B` starts accepting, the OS will select one of the two processes to accept each new connection.  Once `B` is ready:

* `B` sends `SIGUSR1` to `A`
* `B` starts accepting on the shared listener
* `A` receives `SIGUSR1` and calls `StopAccepting()` on the `falcore.Server` instance.  this will
	* stop accepting new connections
	* put falcore into a shutdown mode
	
From here on out, `B` is now your main server process.  Your socket listener was never unavailable.  You have hot restarted!  But what happens to existing connections on `A`?

* any in flight requests will be completed
* keep-alive behavior is disabled so any request that's processed will be the last before the connection is closed in a protocol friendly manner
* any idle connections will be allowed one last request (which will end with Connection:close)
* any idle connections that don't issue requests will be timed out
* once all connections are closed, `ListenAndServe` will return and the application will exit

## Hot restart on Windows

Since Windows doesn't support forkexec, hot restart is not supported on Windows at this time.  If anyone knows how to do this on Windows, we'd love to hear about it.