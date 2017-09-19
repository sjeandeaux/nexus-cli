package repositorymanager

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	lognexuscli "github.com/sjeandeaux/nexus-cli/log"

	"fmt"
	"path/filepath"
)

func init() {
	lognexuscli.Logger.SetOutput(os.Stdout)
}

const (
	expectedPom = `<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd"><modelVersion>4.0.0</modelVersion><groupId>com.jeandeaux</groupId><artifactId>elyne</artifactId><version>0.1.0-SNAPSHOT</version></project>`
	groupID     = "com.jeandeaux"
	artifactID  = "elyne"
	version     = "0.1.0-SNAPSHOT"
)

func TestRepository_generateURL(t *testing.T) {
	repo := NewRepository("http://in.your.anus.fr/third-party", "", "")
	actual := repo.generateURL(&Artifact{GroupID: groupID, ArtifactID: artifactID, Version: version}, "jar")
	expected := "http://in.your.anus.fr/third-party/com/jeandeaux/elyne/0.1.0-SNAPSHOT/elyne-0.1.0-SNAPSHOT.jar"
	if actual != expected {
		t.Fatal("actual", actual, "expected", expected)
	}
}

type call struct {
	called       bool
	calledSha1   bool
	calledMd5    bool
	expected     string
	expectedSha1 string
	expectedMd5  string
}

func (c *call) allIsCalled() bool {
	return c.called && c.calledMd5 && c.calledSha1
}

func TestRepository_UploadArtifact(t *testing.T) {
	forPom := call{
		called:       false,
		calledSha1:   false,
		calledMd5:    false,
		expected:     expectedPom,
		expectedSha1: "1f396c7604363c787362e5916005a0cad72701c0",
		expectedMd5:  "649de9004a8b0e95a7ed1592bcf1ba8c",
	}

	forFile := call{
		called:       false,
		calledSha1:   false,
		calledMd5:    false,
		expected:     "Commodores Nightshift",
		expectedSha1: "3db2e83d419582fbc443067426a6f3cf7b793bcb",
		expectedMd5:  "b32bc418e831a8311b25a455f584879a",
	}

	file, err := ioutil.TempFile(os.TempDir(), ".jar")
	file.WriteString(forFile.expected)
	extension := filepath.Ext(file.Name())

	ts := repositoryManagerImplementation(t, extension, &forPom, &forFile)
	defer ts.Close()
	repo := NewRepository(ts.URL, "bob", "thesponge")

	if err != nil {
		t.Fatal(err)
	}
	a, _ := NewArtifact(groupID, artifactID, version, file.Name())
	err = repo.UploadArtifact(a, "sha1", "md5", "not-found")
	if err != nil {
		t.Fatal(err)
	}

	if !forFile.allIsCalled() {
		t.Errorf("Problem we are waiting more calls %v", forFile)
	}

	if !forPom.allIsCalled() {
		t.Errorf("Problem we are waiting more calls %v", forPom)
	}

}

//repositoryManagerImplementation TODO too big for what is it
func repositoryManagerImplementation(t *testing.T, extension string, forPom *call, forFile *call) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		actual := string(bodyBytes)

		path := "/com/jeandeaux/elyne/0.1.0-SNAPSHOT/elyne-0.1.0-SNAPSHOT"
		switch r.URL.Path {
		case fmt.Sprint(path, ".pom"):
			forPom.called = true
			if actual != forPom.expected {
				if actual != forPom.expected {
					t.Error("actual", actual, "expected", forPom.expected)
				}
			}
		case fmt.Sprint(path, ".pom.sha1"):
			forPom.calledSha1 = true
			if actual != forPom.expectedSha1 {
				if actual != forPom.expectedSha1 {
					t.Error("actual", actual, "expected", forPom.expectedSha1)
				}
			}

		case fmt.Sprint(path, ".pom.md5"):
			forPom.calledMd5 = true
			if actual != forPom.expectedMd5 {
				if actual != forPom.expectedMd5 {
					t.Error("actual", actual, "expected", forPom.expectedMd5)
				}
			}

		case fmt.Sprint(path, extension):
			forFile.called = true
			if actual != forFile.expected {
				if actual != forFile.expected {
					t.Error("actual", actual, "expected", forFile.expected)
				}
			}
		case fmt.Sprint(path, extension, ".sha1"):
			forFile.calledSha1 = true
			if actual != forFile.expectedSha1 {
				if actual != forFile.expectedSha1 {
					t.Error("actual", actual, "expected", forFile.expectedSha1)
				}
			}
		case fmt.Sprint(path, extension, ".md5"):
			forFile.calledMd5 = true
			if actual != forFile.expectedMd5 {
				if actual != forFile.expectedMd5 {
					t.Error("actual", actual, "expected", forFile.expectedMd5)
				}
			}
		}

	}))
	return ts
}

func TestArtifact_writePom(t *testing.T) {
	a := &Artifact{GroupID: groupID, ArtifactID: artifactID, Version: version}
	actualBytes, err := a.writePom()
	if err != nil {
		t.Fatal(err)
	}
	actual := string(actualBytes)

	if actual != expectedPom {
		t.Fatal("actual", actual, "expected", expectedPom)
	}
}

func TestNewArtifact_Ok(t *testing.T) {
	a := &Artifact{GroupID: groupID, ArtifactID: artifactID, Version: version}
	if a.File != "" {
		t.Fatal("actual", a.File, "expected", "nil")
	}
	if a.Pom != nil {
		t.Fatal("actual", a.Pom, "expected", "nil")
	}

	file, err := ioutil.TempFile(os.TempDir(), ".jar")
	if err != nil {
		t.Fatal(err)
	}

	a, err = NewArtifact(groupID, artifactID, version, file.Name())
	if err != nil {
		t.Fatal(err)
	}

	actual := string(a.Pom)

	if actual != expectedPom {
		t.Fatal("actual", actual, "expected", expectedPom)
	}
}

func TestNewArtifact_Ko_Because_Not_Found(t *testing.T) {
	_, err := NewArtifact(groupID, artifactID, version, "<not found>")
	if err == nil {
		t.Fatal("I want a error")
	}
}

func TestNewArtifact_Ko_Because_Name_Empty(t *testing.T) {
	_, err := NewArtifact(groupID, artifactID, version, "")
	if err == nil {
		t.Fatal("I want a error")
	}
	actual := err.Error()
	expected := "You must specify a file"
	if actual != expected {
		t.Fatal("actual", actual, "expected", expected)
	}
}
