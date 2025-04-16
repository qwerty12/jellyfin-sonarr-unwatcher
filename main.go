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
	// TODO: if binding to 0.0.0.0, https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
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

	if j.Item.UserData != nil && j.Item.UserData.IsFavorite != nil && *j.Item.UserData.IsFavorite {
		return
	}

	if j.Series != nil && j.Series.UserData != nil && j.Series.UserData.IsFavorite != nil && *j.Series.UserData.IsFavorite { // oh Go...
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

			if i >= 64 || (1<<i) > uint64(math.MaxInt64)/uint64(time.Second) {
				time.Sleep(time.Duration(math.MaxInt64))
			} else {
				time.Sleep(time.Duration(1<<i) * time.Second)
			}
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
		IdleTimeout:                  time.Minute,
		ReadTimeout:                  30 * time.Second,
		ReadHeaderTimeout:            10 * time.Second,
		WriteTimeout:                 10 * time.Second,
	}
	log.Fatal(s.ListenAndServe())
}
