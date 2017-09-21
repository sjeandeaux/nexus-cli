FROM golang:1.9

WORKDIR /go/src/github.com/sjeandeaux/nexus-cli
COPY . .

RUN go-wrapper download
RUN go-wrapper install -ldflags " \
 -X github.com/sjeandeaux/nexus-cli/information.Version=$(cat VERSION) \
 -X github.com/sjeandeaux/nexus-cli/information.BuildTime=$(date +"%Y-%m-%dT%H:%M:%S") \
 -X github.com/sjeandeaux/nexus-cli/information.GitCommit=$(git rev-parse --short HEAD) \
 -X github.com/sjeandeaux/nexus-cli/information.GitDescribe=$(git describe --tags --always) \
 -X github.com/sjeandeaux/nexus-cli/information.GitDirty=$(test -n "`git status --porcelain`" && echo "+CHANGES" || true)"

ENTRYPOINT ["go-wrapper", "run"]