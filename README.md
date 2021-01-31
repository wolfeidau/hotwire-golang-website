# hotwire-golang-website

This project provides some working examples using [Go](https://golang.org) the [hotwire/turbo](https://turbo.hotwire.dev/) library published by [basecamp](https://basecamp.com/).

# Overview

This service illustrates how to use turbo to enable updates to a website using primarily service side code.

1. Uses [html/template](https://golang.org/pkg/html/template/) for [views](views).
2. Uses [echo](https://echo.labstack.com/) library to simplify routing.
3. Uses [Go 1.16](https://tip.golang.org/doc/go1.16) [Embedded Files](https://tip.golang.org/doc/go1.16#library-embed) to simplify adding templates to binary.
4. Uses a CDN to host all css / JS libraries [base.html](views/layout/base.html).

**Note:** As mentioned this project requires **Go 1.16** which is currently in `rc1`.

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
