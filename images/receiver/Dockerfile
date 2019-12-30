FROM docker.io/library/golang:1.13-alpine as builder

RUN apk add --update git alpine-sdk
RUN mkdir -p /code
ADD . /code/
WORKDIR /code
RUN make binaries

FROM scratch
COPY --from=builder /code/build/wathola-receiver /wathola-receiver
WORKDIR /
CMD ["/wathola-receiver"]
