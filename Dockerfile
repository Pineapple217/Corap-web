# Building the binary of the App
FROM golang:1.20 AS build

WORKDIR /go/src/Corap-web

# Install tailwindcss cli
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64
RUN chmod +x tailwindcss-linux-x64
RUN mv tailwindcss-linux-x64 tailwindcss

COPY . .

# Downloads all the dependencies in advance (could be left out, but it's more clear this way)
RUN go mod download
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

# Builds the application as a staticly linked one, to allow it to run on alpine
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app .


RUN ./tailwindcss -o ./tlw.css --minify


# Moving the binary to the 'final Image' to make it smaller
FROM alpine:latest as release

WORKDIR /app

# Create the `public` dir and copy all the assets into it
RUN mkdir ./static
COPY ./static ./static

COPY --from=build /go/src/Corap-web/app .
COPY --from=build /go/src/Corap-web/tlw.css ./static/public/css/tlw.css

RUN apk -U upgrade \
    && apk add --no-cache dumb-init ca-certificates \
    && chmod +x /app/app

EXPOSE 3000

ENTRYPOINT ["/usr/bin/dumb-init", "--"]