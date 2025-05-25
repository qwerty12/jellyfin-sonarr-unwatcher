//go:generate go tool go-winres make

package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

const PATH_JELLYFIN = "/jellyfin"

func jellyfinHandler(_ http.ResponseWriter, r *http.Request) {
	var j JellyfinPayload
	if !readJellyfinWebhookPayload(r, &j) || j.Item.Id == nil {
		return
	}

	if _, loaded := alreadyUnmonitoredCache.LoadOrCompute(*j.Item.Id, func() (newValue int64, cancel bool) {
		return time.Now().Unix(), false
	}); loaded {
		return
	}

	if !isInSonarrFolder(j.Item.Path) {
		return
	}

	unmonitorEpisode(j.Item, j.Series, nil)
}

func main() {
	sonarrInit()
	go pollInitialRootFolders()

	jellyfinPort := os.Getenv("JELLYFIN_PORT")
	if jellyfinPort == "" {
		jellyfinPort = "9898"
	}
	addr := "127.0.0.1:" + jellyfinPort
	log.Print("Starting unmonitoring for Jellyfin on http://", addr, PATH_JELLYFIN)

	mux := http.NewServeMux()
	mux.HandleFunc("POST "+PATH_JELLYFIN, jellyfinHandler)
	enablePrefetcharr(mux)

	go func() {
		for t := range time.Tick(2 * 24 * time.Hour) {
			nowStart := t.Unix()

			alreadyUnmonitoredCache.Range(func(key string, value int64) bool {
				if nowStart > value {
					alreadyUnmonitoredCache.Delete(key)
				}
				return true
			})

			if alreadyPrefetchedCache != nil {
				alreadyPrefetchedCache.Range(func(key alreadySeenSeason, value int64) bool {
					if nowStart > value {
						alreadyPrefetchedCache.Delete(key)
					}
					return true
				})
			}
		}
	}()

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
