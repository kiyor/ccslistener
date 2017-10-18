FROM golang as builder
COPY . /go/src/github.com/kiyor/ccslistener
RUN cd /go/src/github.com/kiyor/ccslistener && \
    go get && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ccslistener .

FROM docker:dind
WORKDIR /usr/local/bin
COPY --from=builder /go/src/github.com/kiyor/ccslistener/ccslistener .
COPY run.sh .
EXPOSE 8886
VOLUME ["/root/.docker"]

ENTRYPOINT ["./run.sh"]
