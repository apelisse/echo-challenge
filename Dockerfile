FROM golang:1.8 as build

WORKDIR /go/src/echo
COPY . .

RUN go-wrapper download
RUN go-wrapper install

FROM gcr.io/distroless/base
COPY --from=build /go/bin/echo /
ENTRYPOINT ["/echo"]
