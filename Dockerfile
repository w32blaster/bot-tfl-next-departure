#
# Phase 1: build the application
#
FROM w32blaster/go-govendor AS builder

# make right Go project structure
RUN mkdir -p /go/src/github.com/w32blaster/bot-tfl-next-departure/vendor && \
    export GOPATH && \
    GOPATH="/go" 

# copy sources (please refer to .dockerignore file to see what is ignored)
ADD . /go/src/github.com/w32blaster/bot-tfl-next-departure/

RUN cd /go/src/github.com/w32blaster/bot-tfl-next-departure && \
    #
    # download fresh dependencies
    govendor fetch -v +out && \
    #
    # and compile our application
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bot .

#
# Phase 2: prepare the runtime container, ready for production
#
FROM scratch
COPY --from=builder /go/src/github.com/w32blaster/bot-tfl-next-departure/bot /bot

VOLUME "/storage"
CMD ["/bot"]
