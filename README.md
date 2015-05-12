# go-hummingbird #

go-hummingbird is a Go library for accessing the [Hummingbird.me API](https://github.com/hummingbird-me/hummingbird/wiki/API-v1-Methods). Currently only v1 of the API is supported (v2 is still under development).

[![GitHub license](https://img.shields.io/github/license/mashape/apistatus.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/nstratos/go-hummingbird/hb?status.svg)](https://godoc.org/github.com/nstratos/go-hummingbird/hb)
[![Coverage Status](https://coveralls.io/repos/nstratos/go-hummingbird/badge.svg?branch=master)](https://coveralls.io/r/nstratos/go-hummingbird?branch=master)
[![Build Status](https://drone.io/github.com/nstratos/go-hummingbird/status.png)](https://drone.io/github.com/nstratos/go-hummingbird/latest)

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
