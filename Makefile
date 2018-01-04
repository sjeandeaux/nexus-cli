
GOOS?=$(shell uname -s | tr '[:upper:]' '[:lower:]')
GOARCH?=amd64

test:
	go test ./...
nexus:
	docker run -d -p 8081:8081 --name nexus sonatype/nexus3

build: test
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags " \
	     -X github.com/sjeandeaux/nexus-cli/information.Version=$(shell cat VERSION) \
	     -X github.com/sjeandeaux/nexus-cli/information.BuildTime=$(shell date +"%Y-%m-%dT%H:%M:%S") \
	     -X github.com/sjeandeaux/nexus-cli/information.GitCommit=$(shell git rev-parse --short HEAD) \
	     -X github.com/sjeandeaux/nexus-cli/information.GitDescribe=$(shell git describe --tags --always) \
	     -X github.com/sjeandeaux/nexus-cli/information.GitDirty=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)" \
	     -o dist/$(GOOS)-$(GOARCH)-nexus-cli

upload:
	touch upload.jar
	go run main.go -repo=http://localhost:8081/repository/maven-releases \
                              -user=admin \
                              -password=admin123 \
                              -action PUT \
                              -file=upload.jar \
                              -groupID=com.jeandeaux \
                              -artifactID=elyne \
                              -version=0.1.0 \
                              -hash md5 \
                              -hash sha1

delete:
	go run main.go -repo=http://localhost:8081/repository/maven-releases \
                              -user=admin \
                              -password=admin123 \
															-file=upload.jar \
                              -action DELETE \
                              -groupID=com.jeandeaux \
                              -artifactID=elyne \
                              -version=0.1.0 \
															-hash md5 \
                              -hash sha1
