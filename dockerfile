FROM golang:1.16-alpine AS build
RUN apk add -U --no-cache ca-certificates
WORKDIR /src
COPY . .
ENV CGO_ENABLED=0 
RUN go build -o /out/gcp-logger .
FROM scratch AS bin
COPY --from=build /out/gcp-logger /gcp-logger
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/gcp-logger"]
