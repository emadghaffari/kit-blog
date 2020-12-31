FROM golang

RUN mkdir -p /go/src/github.com/emadghaffari/kit-blog

ADD . /go/src/github.com/emadghaffari/kit-blog

RUN go get  -t -v ./...
RUN go get  github.com/canthefason/go-watcher
RUN go install github.com/canthefason/go-watcher/cmd/watcher

ENTRYPOINT  watcher -run github.com/emadghaffari/kit-blog/posts/cmd  -watch github.com/emadghaffari/kit-blog/posts
