# nexus-cli

## upload

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

```

## TODOs

- [ ] interface for artifact
- [ ] interface for repository
- [ ] integration tests
