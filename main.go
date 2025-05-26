//go:generate go tool go-winres make

package main

import (
	"github.com/llxisdsh/pb"
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

	id := *j.Item.Id
	if alreadyUnmonitoredCache.HasKey(id) {
		return
	}
	alreadyUnmonitoredCache.Store(id, time.Now().Unix())

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
		var tickerUnmonitored, tickerPrefetched <-chan time.Time
		tickerUnmonitored = time.Tick(2 * time.Hour)
		if alreadyPrefetchedCache != nil {
			tickerPrefetched = time.Tick(2 * 24 * time.Hour)
		}

		for {
			select {
			case t := <-tickerUnmonitored:
				nowStart := t.Unix()
				//log.Printf("alreadyUnmonitoredCache %s", alreadyUnmonitoredCache.Stats().ToString())
				alreadyUnmonitoredCache.RangeEntry(func(e *pb.EntryOf[string, int64]) bool {
					if nowStart > e.Value {
						alreadyUnmonitoredCache.Delete(e.Key)
					}
					return true
				})
			case t := <-tickerPrefetched:
				nowStart := t.Unix()
				//log.Printf("alreadyPrefetchedCache %s", alreadyPrefetchedCache.Stats().ToString())
				alreadyPrefetchedCache.RangeEntry(func(e *pb.EntryOf[alreadySeenSeason, int64]) bool {
					if nowStart > e.Value {
						alreadyPrefetchedCache.Delete(e.Key)
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
