package version

import (
	_ "embed"
	"strings"
)

//go:embed version.txt
var VersionData string

type VersionInfo struct {
	Tag		  string
	Revision  string
	Dirty	  bool
}

func Version() VersionInfo {
	v := strings.Split(VersionData, "|")
	tag := ""
	rev := ""
	dirty := false
	if(len(v) > 1) {
		tag = v[0]
	}
	if(len(v) > 2) {
		rev = v[1]
	}
	if(len(v) >= 3) {
		dirty = v[2] == "true"
	}
	return VersionInfo{Tag: tag, Revision: rev, Dirty: dirty }
}