//go:generate go tool -modfile=go.tool.mod go-winres make

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/llxisdsh/pb"
)

const PATH_JELLYFIN = "/jellyfin"

var alreadyUnmonitoredCache *pb.MapOf[string, int64]

func jellyfinHandler(_ http.ResponseWriter, r *http.Request) {
	var j JellyfinPayload
	if !readJellyfinWebhookPayload(r, &j) || j.Item.Id == nil {
		return
	}

	if _, loaded := alreadyUnmonitoredCache.LoadOrStoreFn(*j.Item.Id, func() int64 {
		return time.Now().Unix()
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
	alreadyUnmonitoredCache = pb.NewMapOf[string, int64](pb.WithPresize(50), pb.WithShrinkEnabled())
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
		tickerUnmonitored = time.Tick(12 * time.Hour)
		if alreadyPrefetchedCache != nil {
			tickerPrefetched = time.Tick(2 * 24 * time.Hour)
		}

		for {
			select {
			case t := <-tickerUnmonitored:
				nowStart := t.Unix()
				//log.Printf("alreadyUnmonitoredCache %s", alreadyUnmonitoredCache.Stats().ToString())
				alreadyUnmonitoredCache.RangeProcessEntry(func(loaded *pb.EntryOf[string, int64]) *pb.EntryOf[string, int64] {
					if nowStart >= loaded.Value {
						return nil
					}
					return loaded
				})
			case t := <-tickerPrefetched:
				nowStart := t.Unix()
				//log.Printf("alreadyPrefetchedCache %s", alreadyPrefetchedCache.Stats().ToString())
				alreadyPrefetchedCache.RangeProcessEntry(func(loaded *pb.EntryOf[alreadySeenSeason, int64]) *pb.EntryOf[alreadySeenSeason, int64] {
					if nowStart >= loaded.Value {
						return nil
					}
					return loaded
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
