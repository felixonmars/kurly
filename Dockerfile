
FROM golang:1.8 as kurly

RUN go get github.com/davidjpeacock/cli/...
RUN go get github.com/alsm/ioprogress/...
RUN go get github.com/aki237/nscjar/...

COPY . /go/src/github.com/davidjpeacock/kurly

WORKDIR /go/src/github.com/davidjpeacock/kurly

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o kurly *.go \
    && mv kurly /usr/local/bin/

# ==================================================================================

FROM scratch

COPY --from=kurly /usr/local/bin/kurly /kurly

ENTRYPOINT ["/kurly"]
