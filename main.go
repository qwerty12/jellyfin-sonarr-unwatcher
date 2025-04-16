package main

import (
	"jellyfin-sonarr-unwatcher/extmodels/jellygen"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

const PATH_JELLYFIN = "/jellyfin"

var sonarrRootFolders []string
var sonarrRootFoldersSet = false // not atomic.Bool

func jellyfinHandler(_ http.ResponseWriter, r *http.Request) {
	j, err := DecodeJellyfinPayload(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	if j.Item == nil || j.Item.ProviderIds == nil || j.Item.Type == nil {
		return
	}

	if *j.Item.Type != jellygen.BaseItemKindEpisode {
		return
	}

	if j.Item.UserData != nil && j.Item.UserData.IsFavorite != nil && *j.Item.UserData.IsFavorite { // oh Go
		return
	}

	if j.Series != nil && j.Series.UserData != nil && j.Series.UserData.IsFavorite != nil && *j.Series.UserData.IsFavorite {
		return
	}

	if sonarrRootFoldersSet && j.Item.Path != nil {
		itemPath := *j.Item.Path
		inSonarrFolder := false
		for _, rootFolder := range sonarrRootFolders {
			if strings.HasPrefix(itemPath, rootFolder) {
				inSonarrFolder = true
				break
			}
		}

		if !inSonarrFolder {
			return
		}
	}

	unmonitorEpisode(j.Item.ProviderIds, j.Series)
}

func main() {
	sonarrInit()
	go func() {
		const retries = 6
		for i := 0; i <= retries; i++ {
			rootFolders := getRootFolders()
			if rootFolders != nil {
				sonarrRootFolders = rootFolders
				sonarrRootFoldersSet = true
				return
			}

			exp := i
			if i > math.MaxUint32 {
				exp = math.MaxUint32
			}
			var secs uint64
			if exp >= 64 {
				secs = math.MaxUint64
			} else {
				secs = 1 << exp
			}

			time.Sleep(time.Duration(secs) * time.Second)
		}
	}()

	var jellyfinPort = os.Getenv("JELLYFIN_PORT")
	if jellyfinPort == "" {
		jellyfinPort = "9898"
	}
	var addr = "127.0.0.1:" + jellyfinPort
	log.Print("Starting unmonitoring for Jellyfin on http://", addr, PATH_JELLYFIN)

	mux := http.NewServeMux()
	mux.HandleFunc(http.MethodPost+" "+PATH_JELLYFIN, jellyfinHandler)
	s := &http.Server{
		Addr:                         addr,
		Handler:                      mux,
		DisableGeneralOptionsHandler: true,
	}
	log.Fatal(s.ListenAndServe())
}
