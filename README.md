# hotwire-golang-website

This project provides some working examples using [Go](https://golang.org) the [hotwire/turbo](https://turbo.hotwire.dev/) library published by [basecamp](https://basecamp.com/). This is based on a great post about [Hotwire: HTML Over The Wire](https://delitescere.medium.com/hotwire-html-over-the-wire-2c733487268c) by [@delitescere](https://twitter.com/delitescere).

# Overview

This service illustrates how to use turbo to enable updates to a website using primarily service side code.

1. Uses [html/template](https://golang.org/pkg/html/template/) for [views](views).
2. Uses [echo](https://echo.labstack.com/) library to simplify routing.
3. Uses [Go 1.16](https://tip.golang.org/doc/go1.16) [Embedded Files](https://tip.golang.org/doc/go1.16#library-embed) to simplify adding templates to binary.
4. Uses a CDN to host all css / JS libraries [base.html](views/layouts/base.html).

**Note:** As mentioned this project requires **Go 1.16** which is currently in `rc1`.

# Hotwire Turbo

In this site I have implemented:

1. [Turbo Drive](https://turbo.hotwire.dev/handbook/drive)
2. [Turbo Frames](https://turbo.hotwire.dev/handbook/frames)
3. [Turbo Streams](https://turbo.hotwire.dev/handbook/streams) with [Server Sent Events (SSE)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events) and [WebSockets](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket)

Most of the server side logic is in [hotwire.go](internal/server/hotwire.go).

# Running

To get this project running you need to setup some certificates, in my case I use mkcert and there is a target in the makefile.

```
make certs
```

To start the service just run.

```
make start
```

For development you can use the following command, this will check for code updates and restart the service.

```
make watch
```

The service should be listening on https://hotwire.localhost:9443/

# License

This application is released under Apache 2.0 license and is copyright [Mark Wolfe](https://www.wolfe.id.au).
