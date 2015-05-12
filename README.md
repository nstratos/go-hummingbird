# go-hummingbird #

go-hummingbird is a Go library for accessing the [Hummingbird.me API](https://github.com/hummingbird-me/hummingbird/wiki/API-v1-Methods). Currently only v1 of the API is supported (v2 is still under development).

## Installation ##

This package can be installed using:

    go get github.com/nstratos/go-hummingbird/hb

## Usage ##

```go
	import "github.com/nstratos/go-hummingbird/hb"
```

Construct a new client, then use one of the client's services to access the
different Hummingbird API methods. For example, to get the currently watching
anime entries that are contained in the library of the user "cybrox":

```go
	c := hb.NewClient(nil)

	entries, _, err := c.User.Library("cybrox", hb.StatusCurrentlyWatching)
	// handle err

	// do something with entries
```

See more [examples](https://godoc.org/github.com/nstratos/go-hummingbird/hb#pkg-examples).

## License ##
[MIT](LICENSE)