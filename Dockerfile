FROM golang:alpine

COPY ./configs ./configs
COPY ./internal/pkg/templates ./templates
COPY bin/main /main

CMD ["/main"]