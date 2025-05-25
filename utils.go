package main

import (
	"encoding/json"
	"jellyfin-sonarr-unwatcher/internal/jellygen"
	"log"
	"net/http"
	"strconv"
)

func readJellyfinWebhookPayload(r *http.Request, j *JellyfinPayload) bool {
	// TODO: https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body if binding to 0.0.0.0
	if err := json.NewDecoder(r.Body).Decode(j); err != nil {
		log.Print(err)
		return false
	}

	if j.Item == nil || j.Item.Type == nil || *j.Item.Type != jellygen.BaseItemKindEpisode {
		return false
	}

	return true
}

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
