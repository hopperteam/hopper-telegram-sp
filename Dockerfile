FROM golang:alpine as builder
WORKDIR /build
COPY . /build
RUN go build

FROM alpine as runner
EXPOSE 80
COPY --from=builder /build/hopper-telegram-sp /hopper-telegram-sp
COPY ./res /res

ENTRYPOINT ["/hopper-telegram-sp"]
