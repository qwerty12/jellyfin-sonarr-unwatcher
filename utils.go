package main

import (
	"jellyfin-sonarr-unwatcher/internal/jellygen"
	"strconv"
)

func getSeriesIdentifiersFromJfSeries(series *jellygen.BaseItemDto) (seriesTvdbId string, seriesTitle string) {
	if series.ProviderIds != nil {
		seriesTvdbId = (*series.ProviderIds)["Tvdb"]
	}

	if series.Name != nil {
		seriesTitle = *series.Name
	}
	if seriesTitle == "" && series.OriginalTitle != nil {
		seriesTitle = *series.OriginalTitle
	}

	return
}

func atoi32(s string) int32 {
	if i, err := strconv.ParseInt(s, 10, 32); err == nil {
		return int32(i)
	}
	return 0
}

func ptr[T any](val T) *T {
	return &val
}
