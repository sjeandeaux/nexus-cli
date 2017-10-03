FROM golang:1.9
WORKDIR /go/src/github.com/sjeandeaux/nexus-cli
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags " \
   -X github.com/sjeandeaux/nexus-cli/information.Version=$(cat VERSION) \
   -X github.com/sjeandeaux/nexus-cli/information.BuildTime=$(date +"%Y-%m-%dT%H:%M:%S") \
   -X github.com/sjeandeaux/nexus-cli/information.GitCommit=$(git rev-parse --short HEAD) \
   -X github.com/sjeandeaux/nexus-cli/information.GitDescribe=$(git describe --tags --always) \
   -X github.com/sjeandeaux/nexus-cli/information.GitDirty=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)" -a -installsuffix cgo -o nexus-cli .


FROM scratch
COPY --from=0 /go/src/github.com/sjeandeaux/nexus-cli/nexus-cli /
ENTRYPOINT ["/nexus-cli"] 
