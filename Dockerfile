FROM golang:1.19-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/rarimo/kyc-service
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/kyc-service /go/src/github.com/rarimo/kyc-service-legacy


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/kyc-service /usr/local/bin/kyc-service
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["kyc-service"]
