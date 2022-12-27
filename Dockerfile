# based on https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324
############################
# STEP 1 build executable binary
############################
FROM --platform=$BUILDPLATFORM golang:1.19-alpine AS builder
ARG USER=appuser
ARG UID=10001
ARG GO_FILES=.
ARG GO_MAIN=main.go
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
    --uid "${UID}" \
    "${USER}"

# Fetch dependencies.
COPY ${GO_FILES}/go.mod .
COPY ${GO_FILES}/go.sum .
RUN go mod download
# Using go get.
RUN go get -v all
# copy code files
COPY ${GO_FILES} .
# Build the binary.
# TARGETOS and TARGETARCH are builtin vars
ARG TARGETOS
ARG TARGETARCH
ARG CGO_ENABLED=0
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} CGO_ENABLED=${CGO_ENABLED} go build -ldflags="-w -s" -o /app/main.bin ${GO_MAIN}

# create tmp folder for use in scratch
RUN mkdir /my_tmp
RUN chown -R ${UID}:${UID} /my_tmp


############################
# STEP 2 build a small image
############################
FROM scratch
ARG UID=10001
ARG PROJECT_NAME=go-static
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# create tmp directory
COPY --from=builder --chown=${UID}:${UID} /my_tmp /tmp
# Copy static executable.
COPY --from=builder --chown=${UID}:${UID} /app/main.bin /main

ENV UID=$UID
# Use an unprivileged user.
USER ${UID}:${UID}


# Run the binary.
ENTRYPOINT ["/main"]
