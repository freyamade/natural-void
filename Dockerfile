# Golang builder
FROM golang:1.10-alpine as go-build
ENV GOPATH /go/
WORKDIR /go/src/github.com/crnbrdrck/natural-void
RUN apk add git
# Copy in all the go stuff
COPY . .
# Install dependencies
RUN go get -d ./...
# Build the script and install it in /go/bin/
RUN go install
# Output: /go/bin/natural-void/

# npm stuff builder
FROM node:10.11-alpine as node-build
WORKDIR /natural-void/
# Copy in all of the node related stuff
COPY . .
# Install node packages
RUN npm ci
# Compile typescript and SCSS
RUN npm run build
# Output: /natural-void/static

# Runtime image
FROM alpine:latest
LABEL maintainer=crnbrdrck
WORKDIR /natural-void/
# Copy go binary
COPY --from=go-build /go/bin/natural-void ./
# Copy built assets
COPY --from=node-build /natural-void/static/ ./static/
# Copy other necessary stuff
COPY episodes/ .
COPY templates/ ./templates/
# Label volumes
VOLUME /natural-void/episodes
# Expose the port
EXPOSE 3333
# Create the entrypoint
ENTRYPOINT ./natural-void
