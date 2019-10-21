package main

import (
	"log"
	"os"

	"github.com/beamly/go-http-middleware"
	"github.com/valyala/fasthttp"
)

var (
	verbose    = os.Getenv("VERBOSE")
	listenAddr = os.Getenv("LISTEN_ADDR")
)

func main() {
	m := middleware.NewMiddleware(API{verbose: verbose == "true"})

	log.Panic(fasthttp.ListenAndServe(listenAddr, m.ServeFastHTTP))
}
