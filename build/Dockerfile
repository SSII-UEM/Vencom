FROM golang:1.13.8 AS agency_builder

WORKDIR /example
COPY go.mod .
RUN go mod download
COPY cmd cmd
ENV PATH="/example:${PATH}"
RUN cd cmd; CGO_ENABLED=0 GOOS=linux go build -ldflags '-s' -o agency; cp agency /example/

FROM alpine:latest

WORKDIR /root/
#RUN apk add --update netbase ca-certificates
COPY --from=agency_builder /example/agency .
ENV PATH="/root:${PATH}"
EXPOSE 10000
CMD ["./agency"]
