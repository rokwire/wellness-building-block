FROM golang:1.18-buster as builder

ENV CGO_ENABLED=0

RUN mkdir /wellness-app
WORKDIR /wellness-app
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN make

FROM alpine:3.17

#we need timezone database
RUN apk add --no-cache --update tzdata

COPY --from=builder /wellness-app/bin/wellness /
COPY --from=builder /wellness-app/docs/swagger.yaml /docs/swagger.yaml

COPY --from=builder /wellness-app/driver/web/authorization_model.conf /driver/web/authorization_model.conf
COPY --from=builder /wellness-app/driver/web/authorization_policy.csv /driver/web/authorization_policy.csv

ENTRYPOINT ["/wellness"]
