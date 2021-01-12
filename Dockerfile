FROM alpine:3.11

RUN apk add --no-cache ca-certificates && apk update
COPY anodot-kube-events /go/bin/anodot-kube-events

EXPOSE 8080
ENTRYPOINT ["/go/bin/anodot-kube-events"]