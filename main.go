//go:generate go tool go-winres make

package main

import (
	"encoding/json"
	"jellyfin-sonarr-unwatcher/internal/jellygen"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"
)

const PATH_JELLYFIN = "/jellyfin"

var remainingEpisodes int32

var (
	sonarrRootFolders    []string // cached for the process's lifetime
	sonarrRootFoldersSet = false  // not atomic.Bool
)

func isInSonarrFolder(Path *string) bool {
	if sonarrRootFoldersSet && Path != nil {
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

func readJellyfinWebhookPayload(r *http.Request, j *JellyfinPayload) bool {
	// TODO: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body if binding to 0.0.0.0
	err := json.NewDecoder(r.Body).Decode(j)
	if err != nil {
		log.Print(err)
		return false
	}

	if j.Item == nil || j.Item.Type == nil {
		return false
	}

	if *j.Item.Type != jellygen.BaseItemKindEpisode {
		return false
	}

	if !isInSonarrFolder(j.Item.Path) {
		return false
	}

	return true
}

func jellyfinHandler(_ http.ResponseWriter, r *http.Request) {
	var j JellyfinPayload
	if !readJellyfinWebhookPayload(r, &j) {
		return
	}

	unmonitorEpisode(j.Item, j.Series, nil)
}

func prefetcharrHandler(_ http.ResponseWriter, r *http.Request) {
	var j JellyfinPayload
	if !readJellyfinWebhookPayload(r, &j) {
		return
	}

	searchNext(j.Item, j.Series)
}

func main() {
	sonarrInit()
	go func() {
		const retries = 6
		for i := 0; i <= retries; i++ {
			rootFolders := getRootFolders()
			if rootFolders != nil {
				if len(rootFolders) > 0 {
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
	}()

	jellyfinPort := os.Getenv("JELLYFIN_PORT")
	if jellyfinPort == "" {
		jellyfinPort = "9898"
	}
	addr := "127.0.0.1:" + jellyfinPort
	log.Print("Starting unmonitoring for Jellyfin on http://", addr, PATH_JELLYFIN)

	mux := http.NewServeMux()
	mux.HandleFunc("POST "+PATH_JELLYFIN, jellyfinHandler)

	remainingEpisodes = atoi32(os.Getenv("REMAINING_EPISODES"))
	if remainingEpisodes > 0 {
		log.Print("Enabling /prefetcharr endpoint, remaining episodes threshold: ", remainingEpisodes)
		mux.HandleFunc("POST /prefetcharr", prefetcharrHandler)
	}

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
