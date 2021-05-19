# Go-Static

Just another tiny, speedy webserver
Host your static frontend with speed.

## But Why?

This project provides a container that can be used to host static files like an angular, react or vue frontend very easily.


## But How?

There are basically two ways to use this image.
Either you build your own image (preferred) or you mount the static files into a running container.

### Build your own Image

``` Dockerfile
FROM fluktuid/go-static
LABEL maintainer="<yourmail>"

COPY ./your-static-files /static
```

### Mount your Files

``` sh
docker run -d --name go-static -v $(pwd)/your-static-files:/static:ro fluktuid/go-static
```


## But how fast

This Project ist based on [fasthttp](tttps://github.com/valyala/fasthttp) so check [their benchmarks](https://github.com/valyala/fasthttp#http-client-comparison-with-nethttp) for further information.
