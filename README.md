# hotwire-golang-website

This project provides some working examples using [Go](https://golang.org) the [hotwire/turbo](https://turbo.hotwire.dev/) library published by [basecamp](https://basecamp.com/).

# Overview

This service illustrates how to use turbo to enable updates to a website using primarily service side code.

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
