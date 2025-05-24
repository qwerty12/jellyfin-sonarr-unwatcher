package main

import (
	"jellyfin-sonarr-unwatcher/internal/jellygen"
	"jellyfin-sonarr-unwatcher/internal/sonarrt"
	"log"
	"strconv"
	"time"
)

func searchNext(item *jellygen.BaseItemDto, series *jellygen.BaseItemDto) {
	if series == nil {
		return
	}

	var seasonNumber, episodeNumber int32
	if item.IndexNumber != nil {
		episodeNumber = *item.IndexNumber
	}
	if item.ParentIndexNumber != nil {
		seasonNumber = *item.ParentIndexNumber
	}
	if seasonNumber == 0 || episodeNumber == 0 {
		return
	}

	seriesTvdbId, seriesTitle := getSeriesIdentifiersFromJfSeries(series)

	var sonarrSeries *sonarrt.SeriesResource
	for i := range findSonarrSeries(seriesTvdbId, seriesTitle) {
		sonarrSeries = i
		break
	}

	if sonarrSeries == nil {
		//log.Printf("Prefetcharr: series '%s' not found in Sonarr", seriesTitle)
		return
	}

	log.Printf("Prefetcharr: %s S%dE%d", seriesTitle, seasonNumber, episodeNumber)

	season := findSeason(sonarrSeries, seasonNumber)
	if season == nil {
		//log.Print("Season not known to Sonarr")
		return
	}

	isPilot := episodeNumber == 1 && seasonNumber == 1
	isOnlyEpisode := season.Statistics.EpisodeFileCount != nil && *season.Statistics.EpisodeFileCount == 1
	episodecount := int32(0)
	if season.Statistics.TotalEpisodeCount != nil {
		episodecount = *season.Statistics.TotalEpisodeCount
	}
	isEndOfSeason := episodeNumber >= episodecount-remainingEpisodes

	if !(isEndOfSeason || (isPilot && isOnlyEpisode)) {
		//log.Print("Ignoring early episode")
		return
	}

	var nextSeason *sonarrt.SeasonResource
	if isPilot && isOnlyEpisode {
		log.Print("Prefetcharr: Stand-alone pilot episode detected, target first season")
		nextSeason = season
	} else if s := findSeason(sonarrSeries, seasonNumber+1); s != nil {
		nextSeason = s
	} else {
		log.Print("Prefetcharr: Next season not known, monitor new seasons instead")
		sonarrSeries.MonitorNewItems = ptr(sonarrt.NewItemMonitorTypesAll)
		sonarrSeries.Monitored = ptr(true)
		_ = putSeries(sonarrSeries)
		return
	}

	if nextSeason.Statistics.TotalEpisodeCount != nil && nextSeason.Statistics.EpisodeFileCount != nil {
		if *nextSeason.Statistics.TotalEpisodeCount > 0 && *nextSeason.Statistics.TotalEpisodeCount == *nextSeason.Statistics.EpisodeFileCount {
			//log.Print("Prefetcharr: Skip already downloaded season ", *nextSeason.SeasonNumber)
			return
		}
	}

	prefetcharr(sonarrSeries, nextSeason, seasonNumber, item, series)
}

func prefetcharr(sonarrSeries *sonarrt.SeriesResource, season *sonarrt.SeasonResource, jellyfinSeasonNumber int32, item *jellygen.BaseItemDto, series *jellygen.BaseItemDto) {
	seasonNumber := *season.SeasonNumber
	log.Print("Prefetcharr: Searching next season ", seasonNumber)

	var err error
	if season.Monitored == nil || sonarrSeries.Monitored == nil || !*season.Monitored || !*sonarrSeries.Monitored {
		season.Monitored = ptr(true)
		sonarrSeries.Monitored = ptr(true)
		err = putSeries(sonarrSeries)
	}

	if err == nil {
		err = sonarrClient.post("command", nil,
			map[string]any{
				"name":         "SeasonSearch",
				"seasonNumber": seasonNumber,
				"seriesId":     sonarrSeries.Id,
			}, nil)

		if seasonNumber == jellyfinSeasonNumber {
			go func() {
				log.Print("Prefetcharr: waiting two minutes to unmonitor pilot episode again")
				time.Sleep(2 * time.Minute)
				unmonitorEpisode(item, series, sonarrSeries)
			}()
		}
	}

	if err != nil {
		log.Print("Prefetcharr: Error monitoring season: ", err)
	}
}

func findSeason(sonarrSeries *sonarrt.SeriesResource, seasonNumber int32) *sonarrt.SeasonResource {
	if sonarrSeries.Seasons != nil {
		for idx := range *sonarrSeries.Seasons {
			s := &(*sonarrSeries.Seasons)[idx]
			if s.SeasonNumber != nil && *s.SeasonNumber == seasonNumber {
				if s.Statistics != nil {
					return s
				}
				break
			}
		}
	}

	return nil
}

func putSeries(sonarrSeries *sonarrt.SeriesResource) error {
	return sonarrClient.put("series/"+strconv.Itoa(int(*sonarrSeries.Id)), nil, sonarrSeries, nil)
}
