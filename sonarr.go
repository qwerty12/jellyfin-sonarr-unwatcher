package main

import (
	"fmt"
	"jellyfin-sonarr-unwatcher/internal/jellygen"
	"jellyfin-sonarr-unwatcher/internal/sonarrt"
	"log"
	"net/url"
	"os"
	"strconv"
)

var sonarrClient *sonarrAPIClient

func sonarrInit() {
	sonarrHost := os.Getenv("SONARR_HOST")
	sonarrApiKey := os.Getenv("SONARR_API_KEY")
	if sonarrHost == "" || sonarrApiKey == "" {
		log.Fatal("SONARR_HOST and/or $SONARR_API_KEY is required to be set for unmonitoring")
	}

	sonarrHostUrl, err := url.Parse(sonarrHost)
	if err != nil {
		log.Fatal("$SONARR_HOST failed to ", err)
	}
	if sonarrHostUrl.Scheme == "" || sonarrHostUrl.Host == "" {
		log.Fatal("Invalid $SONARR_HOST URL ", sonarrHost)
	}
	sonarrClient = newSonarrAPIClient(sonarrHostUrl, sonarrApiKey)

	log.Println("Sonarr:", sonarrHost)
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

func unmonitorEpisode(episodeProviderIds *map[string]string, series *jellygen.BaseItemDto) {
	var episodeTvdbId int32
	if tvdb := (*episodeProviderIds)["Tvdb"]; tvdb != "" {
		if iTvdb, err := strconv.ParseInt(tvdb, 10, 32); err == nil {
			episodeTvdbId = int32(iTvdb)
		}
	}

	if episodeTvdbId == 0 {
		log.Println("episodeTvdbId is zero")
		return
	}

	var seriesTitle string
	var episode *sonarrt.EpisodeResource
	if episodeSonarrId := (*episodeProviderIds)["sonarr"]; episodeSonarrId != "" {
		episode = findEpisodeBySonarrId(episodeSonarrId, episodeTvdbId)
	}

	if series != nil {
		if series.Name != nil {
			seriesTitle = *series.Name
		}
		if seriesTitle == "" && series.OriginalTitle != nil {
			seriesTitle = *series.OriginalTitle
		}

		if episode == nil {
			var seriesTvdbId string
			if series.ProviderIds != nil {
				seriesTvdbId = (*series.ProviderIds)["Tvdb"]
			}
			if seriesTvdbId != "" {
				episode = findEpisodeBySeriesAndEpisodeTvdbIds(seriesTvdbId, episodeTvdbId)
			} else {
				if seriesTitle == "" {
					log.Printf("SeriesTitle is nil (episodeTvdbId %d)", episodeTvdbId)
					return
				}
				episode = findEpisodeByTitleAndTvdbEpisodeId(seriesTitle, episodeTvdbId)
			}
		}
	} else if episode != nil && episode.Series != nil && episode.Series.Title != nil {
		seriesTitle = *episode.Series.Title
	}

	if episode == nil || episode.SeasonNumber == nil || episode.EpisodeNumber == nil {
		log.Printf("Could not find '%s' in Sonarr library", seriesTitle)
		return
	}

	if episode.Monitored == nil || !*episode.Monitored {
		return
	}

	episodeString := fmt.Sprintf("%s - S%02dE%02d", seriesTitle, *episode.SeasonNumber, *episode.EpisodeNumber)

	if err := sonarrClient.put("episode/monitor", nil, &sonarrt.EpisodesMonitoredResource{
		EpisodeIds: &[]int32{*episode.Id},
		Monitored:  ptr(false),
	}, nil); err != nil {
		log.Printf("Failed to unmonitor %s: %v", episodeString, err)
	} else {
		log.Println(episodeString, "unmonitored!")
	}
}

func findEpisodeBySonarrId(episodeSonarrId string, episodeTvdbId int32) *sonarrt.EpisodeResource {
	var ep sonarrt.EpisodeResource
	if err := sonarrClient.get("episode/"+episodeSonarrId, nil, &ep); err == nil {
		if ep.TvdbId != nil && *ep.TvdbId == episodeTvdbId {
			return &ep
		}
	}

	return nil
}

func findEpisodeBySeriesAndEpisodeTvdbIds(seriesTvdbId string, episodeTvdbId int32) *sonarrt.EpisodeResource {
	var seriesList []sonarrt.SeriesResource
	if err := sonarrClient.get("series", url.Values{
		"tvdbId": []string{seriesTvdbId},
	}, &seriesList); err != nil {
		return nil
	}

	for i := range seriesList {
		series := &seriesList[i]
		if series.Id != nil {
			if ep := findEpisodeInSeries(*series.Id, episodeTvdbId); ep != nil {
				return ep
			}
		}
	}

	return nil
}

func findEpisodeByTitleAndTvdbEpisodeId(seriesTitle string, episodeTvdbId int32) *sonarrt.EpisodeResource {
	// https://github.com/Shraymonks/unmonitorr/blob/main/src/sonarr.ts
	// Sonarr has no api for getting an episode by episode tvdbId
	// Go through the following steps to get the matching episode:
	// 1. Get series list
	// 2. Match potential series on title
	// 3. Get episode lists
	// 4. Match episode on tvdbId
	var seriesList []sonarrt.SeriesResource
	if err := sonarrClient.get("series", nil, &seriesList); err != nil {
		log.Printf("Failed to get series list from Sonarr (to find '%s'): %v", seriesTitle, err)
		return nil
	}

	for i := range seriesList {
		series := &seriesList[i]
		if series.Id != nil && series.Title != nil && *series.Title == seriesTitle {
			if ep := findEpisodeInSeries(*series.Id, episodeTvdbId); ep != nil {
				return ep
			}
		}
	}

	return nil
}

func findEpisodeInSeries(seriesId int32, episodeTvdbId int32) *sonarrt.EpisodeResource {
	var episodeList []sonarrt.EpisodeResource
	if err := sonarrClient.get("episode", url.Values{
		"seriesId": []string{strconv.Itoa(int(seriesId))},
	}, &episodeList); err != nil {
		//log.Printf("Failed to get episode list from Sonarr for '%s': %v", seriesTitle, err)
		return nil
	}

	for _, ep := range episodeList {
		if ep.TvdbId != nil && *ep.TvdbId == episodeTvdbId {
			return &ep
		}
	}

	return nil
}
