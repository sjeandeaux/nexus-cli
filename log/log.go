package log

import (
	"io/ioutil"
	"log"
)

// Logger logger of application.
var Logger *log.Logger = log.New(ioutil.Discard, "[☯ nexus-cli ☯] ⇒ ", log.Ldate|log.Ltime|log.LUTC)
