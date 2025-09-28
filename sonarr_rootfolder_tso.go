//go:build amd64 || 386

package main

import (
	"math"
	"strings"
	"time"

	"jellyfin-sonarr-unwatcher/internal/sonarrt"
)

var (
	sonarrRootFolders    []string // cached for the process's lifetime
	sonarrRootFoldersSet = false  // not atomic.Bool
)

func isInSonarrFolder(Path *string) bool {
	if Path != nil && sonarrRootFoldersSet {
		itemPath := *Path
		for _, rootFolder := range sonarrRootFolders {
			if strings.HasPrefix(itemPath, rootFolder) {
				return true
			}
		}
		return false
	}

	return true
}

func pollInitialRootFolders() {
	const retries = 6
	for i := 0; i <= retries; i++ {
		if rootFolders := getRootFolders(); rootFolders != nil {
			if len(rootFolders) != 0 {
				sonarrRootFolders = rootFolders
				sonarrRootFoldersSet = true
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

func getRootFolders() []string {
	var rootFolders []sonarrt.RootFolderResource
	if err := sonarrClient.get("rootfolder", nil, &rootFolders); err == nil {
		paths := make([]string, 0, len(rootFolders))
		for _, folder := range rootFolders {
			if folder.Path != nil && *folder.Path != "" {
				paths = append(paths, *folder.Path)
			}
		}

		return paths
	}

	return nil
}
