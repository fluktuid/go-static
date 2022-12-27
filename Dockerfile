# based on https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
############################
# STEP 1 build executable binary
############################
FROM golang:1.19-alpine AS builder
ARG GO_FILES=.
ARG GO_MAIN=main.go
ARG USER_NAME=go-static
ARG USER_UID=1001
LABEL maintainer="Lukas Paluch <fluktuid@users.noreply.github.com>"
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache ca-certificates && update-ca-certificates

# Set the Current Working Directory inside the container
WORKDIR /app

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${USER_UID}" \
    "${USER_NAME}"

COPY ${GO_FILES} .
RUN echo ${GO_FILES}
# Fetch dependencies.
RUN go mod download
# Using go get.
RUN go get -v all
# Build the binary.
ARG GOOS=linux
ARG GOARCH=amd64
ARG CGO_ENABLED=0
RUN GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=${CGO_ENABLED} go build -ldflags="-w -s" -o /app/main ${GO_MAIN}

# create tmp folder for use in scratch
RUN mkdir /my_tmp
RUN chown -R ${USER_NAME}:${USER_NAME} /my_tmp


############################
# STEP 2 build a small image
############################
FROM scratch
LABEL maintainer="Lukas Paluch <fluktuid@users.noreply.github.com>"
ARG PROJECT_NAME=go-static
# import certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Copy static executable.
COPY --from=builder /app/main /main
# Copy zoneinfo to get correct timezone
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /
ENV ZONEINFO=/zoneinfo.zip
# set zone to Berlin
ENV TZ=Europe/Berlin
# create tmp directory
COPY --from=builder /my_tmp /tmp

# Use an unprivileged user.
COPY --from=builder /etc/passwd /etc/passwd
USER ${USER_NAME}:${USER_NAME}

# set necessary env
ENV ADDR=":8080"
ENV DIR="/static"

# expose port
EXPOSE 8080/TCP

# Run the binary.
ENTRYPOINT [ "/main" ]

