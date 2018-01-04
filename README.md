# nexus-cli [![Build Status](https://travis-ci.org/sjeandeaux/nexus-cli.svg?branch=master)](https://travis-ci.org/sjeandeaux/nexus-cli)

## TODOs

* 1.0.0 refactoring model because it is the godawful mess

## nexus-cli

I share a volume where i have my upload.jar file.

```bash

#run your own nexus
docker run -d -P --name nexus sonatype/nexus3:3.5.1
docker build --tag sjeandeaux/nexus-cli .
docker run --link nexus:nexus -ti -v $(pwd):$(pwd):ro sjeandeaux/nexus-cli \
                              -repo=http://nexus:8081/repository/maven-releases \
                              -user=admin \
                              -password=admin123 \
                              -file=$(pwd)/upload.jar \
                              -groupID=com.jeandeaux \
                              -artifactID=elyne \
                              -version=0.1.0 \
                              -hash md5 \
                              -hash sha1
#or with the binary
nexus-cli -repo=http://nexus:8081/repository/maven-releases \
                              -user=admin \
                              -password=admin123 \
                              -action PUT \
                              -file=upload.jar \
                              -groupID=com.jeandeaux \
                              -artifactID=elyne \
                              -version=0.1.0 \
                              -hash md5 \
                              -hash sha1

nexus-cli -repo=http://nexus:8081/repository/maven-releases \
                              -user=admin \
                              -password=admin123 \
                              -action DELETE \
                              -file=upload.jar \
                              -groupID=com.jeandeaux \
                              -artifactID=elyne \
                              -version=0.1.0 \
                              -hash md5 \
                              -hash sha1
```
