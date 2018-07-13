package gocrawl

import (
	"net/url"
	"strings"
)

func isParent(parent, child *url.URL) bool {
	if parent.Scheme != child.Scheme {
		return true
	}

	if parent.Host != child.Host {
		return true
	}

	if parent.Port() != child.Port() {
		return true
	}

	parentPath := parent.Path
	if !strings.HasSuffix(parentPath, "/") {
		parentPath = parentPath[:strings.LastIndex(parentPath, "/")+1]
	}
	if !strings.HasPrefix(strings.TrimSuffix(child.Path, "/")+"/", parentPath) {
		return true
	}

	return false
}
