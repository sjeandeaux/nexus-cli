
os=linux
arch=amd64

test:
	go test ./... 

build: test
	GOOS=$(os) GOARCH=$(arch) go build -ldflags " \
	     -X github.com/sjeandeaux/nexus-cli/information.Version=$(shell cat VERSION) \
	     -X github.com/sjeandeaux/nexus-cli/information.BuildTime=$(shell date +"%Y-%m-%dT%H:%M:%S") \
	     -X github.com/sjeandeaux/nexus-cli/information.GitCommit=$(shell git rev-parse --short HEAD) \
	     -X github.com/sjeandeaux/nexus-cli/information.GitDescribe=$(shell git describe --tags --always) \
	     -X github.com/sjeandeaux/nexus-cli/information.GitDirty=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)" \
	     -o dist/$(os)_$(arch)_nexus_cli

