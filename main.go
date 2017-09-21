package main

import (
	"flag"

	"fmt"

	"github.com/sjeandeaux/nexus-cli/log"
	"github.com/sjeandeaux/nexus-cli/repositorymanager"

	"os"
	"github.com/sjeandeaux/nexus-cli/information"
)

type enableHash []string

func (i *enableHash) String() string {
	return fmt.Sprintf("%s", *i)
}

func (i *enableHash) Set(value string) error {
	*i = append(*i, value)
	return nil
}

//commandLineArgs all parameters in command line
type commandLineArgs struct {
	//url of repository to upload file
	urlOfRepository string
	user            string
	password        string
	//file to upload
	file string
	//groupID of artifact
	groupID string
	//artifactID of artifact
	artifactID string
	//version of artifact
	version string
	//hashs the hashs which you want to upload
	hash enableHash
}

var commandLine = &commandLineArgs{}

//init configuration
func init() {
	log.Logger.SetOutput(os.Stdout)
	flag.StringVar(&commandLine.urlOfRepository, "repo", "http://localhost/repository/third-party", "url of repository")
	flag.StringVar(&commandLine.user, "user", "", "user for repository")
	flag.StringVar(&commandLine.password, "password", "", "password for repository")
	flag.StringVar(&commandLine.file, "file", "", "your file to upload on repository")
	flag.StringVar(&commandLine.groupID, "groupID", "com.jeandeaux", "groupid of artifact")
	flag.StringVar(&commandLine.artifactID, "artifactID", "elyne", "artifactID of artifact")
	flag.StringVar(&commandLine.version, "version", "0.1.0-SNAPSHOT", "version of artifact")
	flag.Var(&commandLine.hash, "hash", "md5 or/and sha1")
	flag.Parse()
}

//main upload artifact
func main() {
	log.Logger.Println(information.Print())
	repo := repositorymanager.NewRepository(commandLine.urlOfRepository, commandLine.user, commandLine.password)

	artifact, err := repositorymanager.NewArtifact(commandLine.groupID, commandLine.artifactID, commandLine.version, commandLine.file)
	if err != nil {
		log.Logger.Fatal(err)
	}

	if err := repo.UploadArtifact(artifact, commandLine.hash...); err != nil {
		log.Logger.Fatal(err)
	}
}
