# This version should match that in .nvmrc
FROM --platform=linux/amd64 node:15.11.0 AS nodebuilder

WORKDIR /go/src/github.com/CorriganRenard/kratos-selfservice-ui-go

ADD . .

RUN task clean 
RUN task gen_css

ADD . .

FROM --platform=linux/amd64 golang:1.18 AS gobuilder

WORKDIR /go/src/github.com/CorriganRenard/kratos-selfservice-ui-go

ADD go.mod go.mod
ADD go.sum go.sum

ENV GO111MODULE on
ENV GOOS=linux GOARCH=amd64 CGO_ENABLED=0

RUN go mod download

ADD . .

RUN go build -ldflags="-extldflags=-static" -o /usr/bin/kratos-selfservice-ui-go

FROM --platform=linux/amd64 scratch
COPY --from=gobuilder /usr/bin/kratos-selfservice-ui-go /

# Expose the default port that we will be listening to
EXPOSE 4455

ENTRYPOINT ["/kratos-selfservice-ui-go"]