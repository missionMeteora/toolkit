# Closer

# Closer [![GoDoc](https://godoc.org/github.com/missionMeteora/closer?status.svg)](https://godoc.org/github.com/missionMeteora/closer) ![Status](https://img.shields.io/badge/status-beta-yellow.svg)

Closer is a library which handles signal-catching for proper closing of your Go services.

## Usage
``` go
package main

import (
	"fmt"
	"net/http"

	"github.com/missionMeteora/toolkit/closer"
)

func main() {
	s := newSrv()
	go s.listen()
	s.Close(s.c.Wait())
}

func newSrv() *srv {
	return &srv{
		c: closer.New(),
	}
}

type srv struct {
	c *closer.Closer
}

func (s *srv) listen() {
	s.c.Close(http.ListenAndServe(":80", s))
}

func (s *srv) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Serve http here
}

func (s *srv) Close(err error) {
	fmt.Println("Closing service:", err)
	// Close internal services here
}
```