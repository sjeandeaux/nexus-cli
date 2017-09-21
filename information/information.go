package information

import "fmt"

//We use ldflags
var (
	Version     = "No Version Provided"
	GitCommit   = "No GitCommit Provided"
	GitDescribe = "No GitDescribe Provided"
	GitDirty    = "No GitDirty Provided"
	BuildTime   = "No BuildTime Provided"
)


func Print() string {
	return fmt.Sprintf("%q-%q-%q-%q-%q", Version, BuildTime, GitCommit, GitDescribe, GitDirty)
}