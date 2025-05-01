package main

import (
	"fmt"
	"iter"
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

func unmonitorEpisode(episode *jellygen.BaseItemDto, series *jellygen.BaseItemDto, sonarrSeries *sonarrt.SeriesResource) {
	if episode.ProviderIds == nil {
		return
	}

	if episode.UserData != nil && episode.UserData.IsFavorite != nil && *episode.UserData.IsFavorite {
		return
	}

	if series != nil && series.UserData != nil && series.UserData.IsFavorite != nil && *series.UserData.IsFavorite { // oh Go...
		return
	}

	var episodeTvdbId int32
	if tvdb := (*episode.ProviderIds)["Tvdb"]; tvdb != "" {
		episodeTvdbId = atoi32(tvdb)
	}

	if episodeTvdbId == 0 {
		log.Print("TVDB episode ID not found")
		return
	}

	var seriesTitle string
	var sonarrEpisode *sonarrt.EpisodeResource
	if episodeSonarrId := (*episode.ProviderIds)["sonarr"]; episodeSonarrId != "" {
		sonarrEpisode = findEpisodeBySonarrId(episodeSonarrId, episodeTvdbId)
	}

	if series != nil {
		var seriesTvdbId string
		seriesTvdbId, seriesTitle = getSeriesIdentifiersFromJfSeries(series)

		if sonarrEpisode == nil {
			if sonarrSeries == nil {
				sonarrEpisode = findEpisodeByTvdbIdsOrTitle(seriesTvdbId, seriesTitle, episodeTvdbId)
			} else {
				sonarrEpisode = findEpisodeInSeries(*sonarrSeries.Id, episodeTvdbId)
			}
		}
	} else if sonarrEpisode != nil && sonarrEpisode.Series != nil && sonarrEpisode.Series.Title != nil {
		seriesTitle = *sonarrEpisode.Series.Title
	}

	if sonarrEpisode == nil || sonarrEpisode.SeasonNumber == nil || sonarrEpisode.EpisodeNumber == nil {
		log.Printf("Could not find '%s' in Sonarr library", seriesTitle)
		return
	}

	if sonarrEpisode.Monitored == nil || !*sonarrEpisode.Monitored {
		return
	}

	episodeString := fmt.Sprintf("%s - S%02dE%02d", seriesTitle, *sonarrEpisode.SeasonNumber, *sonarrEpisode.EpisodeNumber)

	if err := sonarrClient.put("episode/monitor", nil, &sonarrt.EpisodesMonitoredResource{
		EpisodeIds: &[]int32{*sonarrEpisode.Id},
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

func findEpisodeByTvdbIdsOrTitle(seriesTvdbId string, seriesTitle string, episodeTvdbId int32) *sonarrt.EpisodeResource {
	for series := range findSonarrSeries(seriesTvdbId, seriesTitle) {
		if ep := findEpisodeInSeries(*series.Id, episodeTvdbId); ep != nil {
			return ep
		}
	}

	return nil
}

func findSonarrSeries(seriesTvdbId string, seriesTitle string) iter.Seq[*sonarrt.SeriesResource] {
	return func(yield func(*sonarrt.SeriesResource) bool) {
		var seriesList []sonarrt.SeriesResource

		if seriesTvdbId != "" {
			//log.Printf("Searching by TVDB ID: %s\n", seriesTvdbId)
			if err := sonarrClient.get("series", url.Values{
				"tvdbId": []string{seriesTvdbId},
			}, &seriesList); err == nil {
				for i := range seriesList {
					series := &seriesList[i]
					if series.Id != nil && !yield(series) {
						return
					}
				}
			} /* else {
				log.Printf("Error fetching series by TVDB ID %s: %v\n", seriesTvdbId, err)
			}*/

			return
		}

		if seriesTitle != "" {
			// https://github.com/Shraymonks/unmonitorr/blob/main/src/sonarr.ts
			// Sonarr has no api for getting an episode by episode tvdbId
			// Go through the following steps to get the matching episode:
			// 1. Get series list
			// 2. Match potential series on title
			// 3. Get episode lists
			// 4. Match episode on tvdbId
			//log.Printf("Searching by title: %s\n", seriesTitle)
			if err := sonarrClient.get("series", nil, &seriesList); err == nil {
				for i := range seriesList {
					series := &seriesList[i]
					if series.Id != nil && series.Title != nil && *series.Title == seriesTitle {
						if !yield(series) {
							return
						}
					}
				}
			} /* else {
				log.Printf("Failed to get list of all series from Sonarr (to find '%s' by its title): %v", seriesTitle, err)
			}*/

			return
		}

		//log.Println("No search criteria provided.")
	}
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
