//
// Serves static files from the given directory.
// Exports various stats at /stats .
//
// gratefully copied from https://github.com/valyala/fasthttp/blob/master/examples/fileserver/fileserver.go
package main

import (
	"expvar"
	"flag"
	"log"

	"github.com/kouhin/envflag"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
)

var (
	addr               = flag.String("addr", ":8080", "TCP address to listen to")
	addrTLS            = flag.String("addrTLS", "", "TCP address to listen to TLS (aka SSL or HTTPS) requests. Leave empty for disabling TLS")
	addrStats          = flag.String("addrStats", "", "TCP address to serve stats server on")
	byteRange          = flag.Bool("byteRange", false, "Enables byte range requests if set to true")
	certFile           = flag.String("certFile", "./ssl-cert.pem", "Path to TLS certificate file")
	compress           = flag.Bool("compress", false, "Enables transparent response compression if set to true")
	dir                = flag.String("dir", "/static", "Directory to serve static files from")
	generateIndexPages = flag.Bool("generateIndexPages", true, "Whether to generate directory index pages")
	keyFile            = flag.String("keyFile", "./ssl-cert.key", "Path to TLS key file")
	vhost              = flag.Bool("vhost", false, "Enables virtual hosting by prepending the requested path with the requested hostname")
	stats              = flag.Bool("stats", true, "Enables stats serving")
)

func main() {
	// Parse command-line flags.
	envflag.Parse()

	// Setup FS handler
	fs := &fasthttp.FS{
		Root:               *dir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: *generateIndexPages,
		Compress:           *compress,
		AcceptByteRange:    *byteRange,
	}
	if *vhost {
		fs.PathRewrite = fasthttp.NewVHostPathRewriter(0)
	}
	fsHandler := fs.NewRequestHandler()

	// Start HTTP server.
	if len(*addr) > 0 {
		log.Printf("Starting HTTP server on %q", *addr)
		go func() {
			handler := getRequestHandler(*stats, *addrStats, fsHandler)
			if err := fasthttp.ListenAndServe(*addr, handler); err != nil {
				log.Fatalf("error in ListenAndServe: %s", err)
			}
		}()
	}

	// Start HTTPS server.
	if len(*addrTLS) > 0 {
		log.Printf("Starting HTTPS server on %q", *addrTLS)
		go func() {
			handler := getRequestHandler(*stats, *addrStats, fsHandler)
			if err := fasthttp.ListenAndServeTLS(*addrTLS, *certFile, *keyFile, handler); err != nil {
				log.Fatalf("error in ListenAndServeTLS: %s", err)
			}
		}()
	}

	// Start stats server.
	if len(*addrStats) > 0 {
		log.Printf("Starting HTTPS server on %q", *addrTLS)
		go func() {
			handler := func(ctx *fasthttp.RequestCtx) {
				expvarhandler.ExpvarHandler(ctx)
			}
			if err := fasthttp.ListenAndServeTLS(*addrTLS, *certFile, *keyFile, handler); err != nil {
				log.Fatalf("error in ListenAndServeTLS: %s", err)
			}
		}()
	}

	log.Printf("Serving files from directory %q", *dir)
	if *stats && len(*addrStats) == 0 {
		log.Printf("See stats at http://%s/stats", *addr)
	} else if len(*addrStats) > 0 {
		log.Printf("See stats at http://%s/", *addrStats)
	}

	// Wait forever.
	select {}
}

func getRequestHandler(stats bool, addrStats string, fsHandler fasthttp.RequestHandler) func(*fasthttp.RequestCtx) {
	if stats && len(addrStats) == 0 {
		// Create RequestHandler serving server stats on /stats and files
		// on other requested paths.
		// /stats output may be filtered using regexps. For example:
		//
		//   * /stats?r=fs will show only stats (expvars) containing 'fs'
		//     in their names.
		return func(ctx *fasthttp.RequestCtx) {
			switch string(ctx.Path()) {
			case "/stats":
				expvarhandler.ExpvarHandler(ctx)
			default:
				fsHandler(ctx)
				updateFSCounters(ctx)
			}
		}
	} else {
		return func(ctx *fasthttp.RequestCtx) {
			fsHandler(ctx)
		}
	}
}

func updateFSCounters(ctx *fasthttp.RequestCtx) {
	// Increment the number of fsHandler calls.
	fsCalls.Add(1)

	// Update other stats counters
	resp := &ctx.Response
	switch resp.StatusCode() {
	case fasthttp.StatusOK:
		fsOKResponses.Add(1)
		fsResponseBodyBytes.Add(int64(resp.Header.ContentLength()))
	case fasthttp.StatusNotModified:
		fsNotModifiedResponses.Add(1)
	case fasthttp.StatusNotFound:
		fsNotFoundResponses.Add(1)
	default:
		fsOtherResponses.Add(1)
	}
}

// Various counters - see https://golang.org/pkg/expvar/ for details.
var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)
