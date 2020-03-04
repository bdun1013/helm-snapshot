FROM golang:1.14.0-alpine3.11 as ALPINE-BUILDER
RUN apk --no-cache add --quiet alpine-sdk=0.5-r0
WORKDIR /go/src/github.com/bpdunni/helm-unittest/
COPY . .
RUN install -d /opt && make install HELM_PLUGIN_DIR=/opt

FROM alpine:3.7 as ALPINE
COPY --from=ALPINE-BUILDER /opt /opt
