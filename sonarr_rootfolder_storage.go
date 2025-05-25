package main

import (
	"math"
	"strings"
	"sync/atomic"
	"time"
)

var sonarrRootFolders atomic.Pointer[[]string] // cached for the process's lifetime

func isInSonarrFolder(Path *string) bool {
	if Path != nil {
		if sonarrRootFolders := sonarrRootFolders.Load(); sonarrRootFolders != nil {
			itemPath := *Path
			for _, rootFolder := range *sonarrRootFolders {
				if strings.HasPrefix(itemPath, rootFolder) {
					return true
				}
			}
			return false
		}
	}

	return true
}

func pollInitialRootFolders() {
	const retries = 6
	for i := 0; i <= retries; i++ {
		if rootFolders := getRootFolders(); rootFolders != nil {
			if len(*rootFolders) != 0 {
				sonarrRootFolders.Store(rootFolders)
			}
			return
		}

		if i <= 33 {
			time.Sleep(time.Duration(1<<i) * time.Second)
		} else {
			time.Sleep(time.Duration(math.MaxInt64))
		}
	}
}
