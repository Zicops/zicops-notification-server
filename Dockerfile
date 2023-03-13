FROM golang:1.18 as builder
RUN apt-get update -qq

ARG GO_MODULES_TOKEN
ENV GO111MODULE=on
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o zicops-notification-server .

FROM golang:latest
LABEL maintainer="Puneet Saraswat <puneet.saraswat10074@gmail.com>"

RUN apt-get update -y -qq

COPY --from=builder /go/src/app/zicops-notification-server /usr/bin/
EXPOSE 8094

ENTRYPOINT ["/usr/bin/zicops-notification-server"]
