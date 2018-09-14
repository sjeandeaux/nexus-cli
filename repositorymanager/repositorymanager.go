//Package repositorymanager upload a file in this repository
package repositorymanager

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"hash"
	"io"
	"net/http"
	"os"
	"strings"
	"text/template"

	"net/url"

	"github.com/sjeandeaux/nexus-cli/log"

	"github.com/sjeandeaux/nexus-cli/information"
)

const (
	dot            = "."
	suffixPom      = "pom"
	allReplacement = -1
	slash          = "/"
	dash           = "-"
)

//Repository where we want put the file
type Repository struct {
	url      string
	user     string
	password string
	client   *http.Client
	hash     map[string]*repositoryHash
}

// repositoryHash create hash and has the suffix for the file on repository
type repositoryHash struct {
	suffix string
	//TODO see if we need to a func or variable
	hash func() hash.Hash
}

//NewRepository create a Repository with default client HTTP.
func NewRepository(url, user, password string) *Repository {
	const (
		nameMd5  = "md5"
		nameSha1 = "sha1"
	)

	shaOneAndMdFive := make(map[string]*repositoryHash)

	shaOneAndMdFive[nameMd5] = &repositoryHash{
		suffix: nameMd5,
		hash:   func() hash.Hash { return md5.New() },
	}

	shaOneAndMdFive[nameSha1] = &repositoryHash{
		suffix: nameSha1,
		hash:   func() hash.Hash { return sha1.New() },
	}

	return &Repository{
		url:      url,
		user:     user,
		password: password,
		client:   &http.Client{},
		hash:     shaOneAndMdFive}

}

//DeleteArtifact deletes the artifact
//ar the artifact to delete
//hashs list of hash to delete
func (n *Repository) DeleteArtifact(ar *Artifact, hashs ...string) error {
	pomURL := n.generateURL(ar, suffixPom)
	if err := n.delete(pomURL); err != nil {
		return err
	}

	fileURL := n.generateURL(ar, ar.extension())
	if err := n.delete(fileURL); err != nil {
		return err
	}

	for _, h := range hashs {
		//if we can't delete hash we continue
		n.deleteHash(ar, h)
	}
	return nil
}

//UploadArtifact upload ar on repository TODO goroutine to upload
//ar the artifact to upload
//hashs list of hash to send
func (n *Repository) UploadArtifact(ar *Artifact, hashs ...string) error {
	pomURL := n.generateURL(ar, suffixPom)
	if err := n.upload(pomURL, bytes.NewReader(ar.Pom), ""); err != nil {
		return err
	}

	fOpen, err := os.Open(ar.File)
	if err != nil {
		return err
	}

	fileURL := n.generateURL(ar, ar.extension())
	if err := n.upload(fileURL, fOpen, ar.ContentType); err != nil {
		return err
	}

	for _, h := range hashs {
		if iGetIt := n.hash[h]; iGetIt != nil {
			n.uploadHash(ar, iGetIt)
		} else {
			urlIssue := generateURLIssue(h)
			log.Logger.Printf("%q is not managed by the client %q", h, urlIssue)
		}
	}

	return nil
}

//generateURLIssue generate the URL on github to create issue on hash method
func generateURLIssue(h string) string {
	const (
		title      = "Move your ass"
		urlFormat  = "https://github.com/sjeandeaux/nexus-cli/issues/new?title=%s&body=%s"
		bodyFormat = "Could you add the hash %q lazy man?\n%s"
	)
	escapedTitle := url.QueryEscape(title)
	body := fmt.Sprintf(bodyFormat, h, information.Print())
	escapedBody := url.QueryEscape(body)
	urlIssue := fmt.Sprintf(urlFormat, escapedTitle, escapedBody)
	return urlIssue
}

//generateURL generate the url of ar
// <url>/<groupID/<arID>/<version>/<arID>-<version>.<endOfFile>
func (n *Repository) generateURL(ar *Artifact, endOfFile string) string {
	g := strings.Replace(ar.GroupID, dot, slash, allReplacement)
	nameOfFile := fmt.Sprint(slash, ar.ArtifactID, dash, ar.Version, dot, endOfFile)
	return fmt.Sprint(n.url, slash, g, slash, ar.ArtifactID, slash, ar.Version, nameOfFile)
}

//uploadHash upload the hash
func (n *Repository) uploadHash(ar *Artifact, h *repositoryHash) error {
	p, f, err := generateHash(ar.Pom, ar.File, h.hash)
	if err != nil {
		return err
	}

	hashedPom := n.generateURL(ar, fmt.Sprint(suffixPom, dot, h.suffix))
	if err = n.upload(hashedPom, p, ""); err != nil {
		return err
	}

	hashedFile := n.generateURL(ar, fmt.Sprint(ar.extension(), dot, h.suffix))
	return n.upload(hashedFile, f, "")
}

//deleteHash delete the hash
func (n *Repository) deleteHash(ar *Artifact, h string) error {
	hashedPom := n.generateURL(ar, fmt.Sprint(suffixPom, dot, h))
	if err := n.delete(hashedPom); err != nil {
		return err
	}

	hashedFile := n.generateURL(ar, fmt.Sprint(ar.extension(), dot, h))
	return n.delete(hashedFile)
}

func (n *Repository) upload(url string, data io.Reader, contentType string) error {
	const (
		PUT               = "PUT"
		httpSuccess       = 201
		HeaderContentType = "Content-Type"
	)

	log.Logger.Print(url)
	req, _ := http.NewRequest(PUT, url, data)

	if n.user != "" && n.password != "" {
		req.SetBasicAuth(n.user, n.password)
	}

	if contentType != "" {
		req.Header.Set(HeaderContentType, contentType)
	}
	res, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != httpSuccess {
		return fmt.Errorf(res.Status)
	}
	return nil
}

func (n *Repository) delete(url string) error {
	const httpSuccess = 204
	log.Logger.Print(url)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	if n.user != "" && n.password != "" {
		req.SetBasicAuth(n.user, n.password)
	}

	res, err := n.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != httpSuccess {
		return fmt.Errorf(res.Status)
	}
	return nil
}

//generateHash generate the
func generateHash(pom []byte, file string, h func() hash.Hash) (io.Reader, io.Reader, error) {
	hashedPom, err := hashSum(bytes.NewReader(pom), h())
	if err != nil {
		return nil, nil, err
	}

	f, errOnFile := os.Open(file)
	defer f.Close()
	if errOnFile != nil {
		return nil, nil, errOnFile
	}

	hashedFile, err := hashSum(f, h())
	if err != nil {
		return nil, nil, err
	}

	return hashedPom, hashedFile, nil
}

//generate the hash of io.Reader
func hashSum(data io.Reader, h hash.Hash) (io.Reader, error) {
	if _, err := io.Copy(h, data); err != nil {
		return nil, err
	}
	return strings.NewReader(hex.EncodeToString(h.Sum(nil))), nil
}

//Artifact the artifact
type Artifact struct {
	//GroupID of artifact
	GroupID string
	//ArtifactID of artifact
	ArtifactID string
	//Version of artifact
	Version string
	//ContentType for the PUT
	ContentType string
	//file to upload
	File string
	//pom of this artifact
	Pom []byte
}

//NewArtifact create a artifact with this own pom
func NewArtifact(groupID, artifactID, version, contentType, file string) (*Artifact, error) {

	if file == "" {
		return nil, errors.New("You must specify a file")

	}

	a := &Artifact{
		GroupID:     groupID,
		ArtifactID:  artifactID,
		Version:     version,
		File:        file,
		ContentType: contentType,
	}

	pom, err := a.writePom()
	if err != nil {
		return nil, err
	}
	a.Pom = pom
	return a, nil
}

// extension extension of file
func (artifact *Artifact) extension() string {
	const unknown = ""
	i := strings.Index(artifact.File, ".")

	if i != -1 {
		i = 1 + i
		return artifact.File[i:]
	}
	return unknown
}

// writePom write a wonderful pom
func (artifact *Artifact) writePom() ([]byte, error) {
	const templateName = "pom"
	const templateValue = `<?xml version="1.0" encoding="UTF-8"?><project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd"><modelVersion>4.0.0</modelVersion><groupId>{{.GroupID}}</groupId><artifactId>{{.ArtifactID}}</artifactId><version>{{.Version}}</version></project>`

	pomTemplate, err := template.New(templateName).Parse(templateValue)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = pomTemplate.Execute(&buf, artifact)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
