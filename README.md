# Go-Static

_Just another tiny, speedy webserver_\
Host your static frontend with speed.

(The code is based on the [fasthttp fileserver example](https://github.com/valyala/fasthttp/blob/master/examples/fileserver/fileserver.go).)

## But Why?

This project provides a container that can be used to host static files like an angular, react or vue frontend very easily.


## But How?

There are basically two ways to use this image.
Either you build your own image (preferred) or you mount the static files into a running container.

### Build your own Image

``` Dockerfile
FROM fluktuid/go-static:1.0.1
# or use the github registry: FROM docker.pkg.github.com/fluktuid/go-static/go-static:1.0.1
LABEL maintainer="yourname <yourmail>"

COPY ./your-static-files /static
```

### Mount your Files

``` sh
docker run -d --name go-static -v $(pwd)/your-static-files:/static:ro fluktuid/go-static
```

## I need more configuration

Of course, a more precise configuration of the service is also possible.
There are different values that can be configured via environment variables:


env var              | description                                                                             | default 
-------------------- | --------------------------------------------------------------------------------------- | -------------- |
ADDR                 | TCP address to listen to                                                                | :8080          |
ADDR_TLS             | TCP address to listen to TLS (aka SSL or HTTPS) requests. Leave empty for disabling TLS | ""             |
BYTE_RANGE           | Enables byte range requests if set to true                                              | false          |
CERT_FILE            | Path to TLS certificate file                                                            | ssl-cert.pem   |
COMPRESS             | Enables transparent response compression if set to true                                 | false          |
DIR                  | Directory to serve static files from                                                    | /static        |
GENERATE_INDEX_PAGES | Whether to generate directory index pages                                               | true           |
KEY_FILE             | Path to TLS key file                                                                    | ./ssl-cert.key |
VHOST                | Enables virtual hosting by prepending the requested path with the requested hostname    | false          |

## But how fast

This Project ist based on [fasthttp](tttps://github.com/valyala/fasthttp) so check [their benchmarks](https://github.com/valyala/fasthttp#http-client-comparison-with-nethttp) for further information.
