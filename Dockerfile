FROM golang:1.24-alpine AS builder

ENV CGO_ENABLED=0

RUN mkdir /wellness-app
WORKDIR /wellness-app
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN make

FROM alpine:3.21.3

#we need timezone database + certificates
RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /wellness-app/bin/wellness /
COPY --from=builder /wellness-app/docs/swagger.yaml /docs/swagger.yaml

COPY --from=builder /wellness-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_scope.conf /wellness-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_scope.conf
COPY --from=builder /wellness-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_string.conf /wellness-app/vendor/github.com/rokwire/rokwire-building-block-sdk-go/services/core/auth/authorization/authorization_model_string.conf

ENTRYPOINT ["/wellness"]
