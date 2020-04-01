# Default to Go 1.12
ARG GO_VERSION=1.12
#First stage: build
FROM golang:${GO_VERSION}-alpine AS builder
# from dm-frontend

LABEL version="0.0.1" vendor="Yadro"

# Create the user and run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN apk add --no-cache ca-certificates git

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
#COPY ./go.mod ./go.sum ./
#RUN go mod download

# Import the code from the context.
COPY ./ ./

#Run unit tests
RUN CGO_ENABLED=0 GOOS="$(go env GOOS)" GOARCH="$(go env GOARCH)" go test -mod=vendor ./...

# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 GOOS="$(go env GOOS)" GOARCH="$(go env GOARCH)" go build \
    -mod=vendor \
    -installsuffix 'static' \
    -ldflags "-s -w" \
    -o /app ./cmd/ch-server/

# Build the directory to upload metadata.
RUN mkdir ./uploads
RUN chown -R nobody:nobody ./uploads
RUN chmod -R 777 ./uploads

# Build the directory contains config.
RUN chown -R nobody:nobody ./config
RUN chmod -R 777 ./config

# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the compiled executable from the first stage.
COPY --from=builder /app /app

# Import the directory to upload metadata.
COPY --from=builder --chown=nobody:nobody /src/uploads /uploads

# Import the directory to contains config.
COPY --from=builder --chown=nobody:nobody /src/config/config.yaml /config/

# Declare the port on which the webserver will be exposed.
# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 8080

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/app"]
