# Building the binary of the App
FROM golang:1.20 AS build

WORKDIR /go/src/Corap-web

COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest as release

WORKDIR /app

# Create the `public` dir and copy all the assets into it
RUN mkdir ./static
COPY ./static ./static

RUN mkdir ./views
COPY ./views ./views

COPY --from=build /go/src/Corap-web/app .

RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/app

EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]