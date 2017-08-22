FROM golang:1.6
WORKDIR /go/src/github.com/drone-plugins/drone-webhook
ADD . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /tmp/drone-webhook .

FROM centurylink/ca-certs
COPY --from=0 /tmp/drone-webhook /bin/
ENTRYPOINT ["/bin/drone-webhook"]

