package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jellyfin-sonarr-unwatcher/extmodels/jellygen"
	"jellyfin-sonarr-unwatcher/extmodels/sonarrt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
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
	var episode sonarrt.EpisodeResource
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

		if episode.TvdbId == nil && series.ProviderIds != nil {
			if seriesTvdbId := (*series.ProviderIds)["Tvdb"]; seriesTvdbId != "" {
				episode = findEpisodeBySeriesAndEpisodeTvdbIds(seriesTvdbId, episodeTvdbId)
			}
		}

		if episode.TvdbId == nil {
			if seriesTitle == "" {
				log.Println("SeriesTitle is nil")
				return
			}
			episode = findEpisodeByTitleAndTvdbEpisodeId(seriesTitle, episodeTvdbId)
		}
	}

	if episode.TvdbId == nil || episode.SeasonNumber == nil || episode.EpisodeNumber == nil {
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

func getRootFolders() []string {
	var paths []string
	var rootFolders []sonarrt.RootFolderResource
	if err := sonarrClient.get("rootfolder", nil, &rootFolders); err == nil {
		for _, folder := range rootFolders {
			if folder.Path != nil && *folder.Path != "" {
				paths = append(paths, *folder.Path)
			}
		}
	}

	return paths
}

func findEpisodeBySonarrId(episodeSonarrId string, episodeTvdbId int32) (matchingEpisode sonarrt.EpisodeResource) {
	var ep sonarrt.EpisodeResource
	if err := sonarrClient.get("episode/"+episodeSonarrId, nil, &ep); err == nil {
		if ep.TvdbId != nil && *ep.TvdbId == episodeTvdbId {
			matchingEpisode = ep
		}
	}
	return
}

func findEpisodeBySeriesAndEpisodeTvdbIds(seriesTvdbId string, episodeTvdbId int32) (matchingEpisode sonarrt.EpisodeResource) {
	var seriesList []sonarrt.SeriesResource
	if err := sonarrClient.get("series", map[string]string{
		"tvdbId": seriesTvdbId,
	}, &seriesList); err != nil {
		return
	}

	for i := range seriesList {
		series := &seriesList[i]
		if series.Id != nil {
			var episodeList []sonarrt.EpisodeResource
			if err := sonarrClient.get("episode", map[string]string{
				"seriesId": strconv.Itoa(int(*series.Id)),
			}, &episodeList); err != nil {
				continue
			}

			for _, ep := range episodeList {
				if ep.TvdbId != nil && *ep.TvdbId == episodeTvdbId {
					matchingEpisode = ep
					return
				}
			}
		}
	}

	return
}

func findEpisodeByTitleAndTvdbEpisodeId(seriesTitle string, episodeTvdbId int32) (matchingEpisode sonarrt.EpisodeResource) {
	// https://github.com/Shraymonks/unmonitorr/blob/main/src/sonarr.ts
	// Sonarr has no api for getting an episode by episode tvdbId
	// Go through the following steps to get the matching episode:
	// 1. Get series list
	// 2. Match potential series on title
	// 3. Get episode lists
	// 4. Match episode on tvdbId
	var seriesList []sonarrt.SeriesResource
	if err := sonarrClient.get("series", nil, &seriesList); err != nil {
		log.Printf("Failed to get series list from Sonarr for '%s': %v", seriesTitle, err)
		return
	}

	cleanedTitle := cleanTitle(seriesTitle)
	// "Match potential series on title. Year metadata from Plex is for the episode
	// so cannot be used for series filtering."
	for i := range seriesList {
		series := &seriesList[i]
		if series.Id != nil && series.Title != nil && cleanTitle(*series.Title) == cleanedTitle {
			var episodeList []sonarrt.EpisodeResource
			if err := sonarrClient.get("episode", map[string]string{
				"seriesId": strconv.Itoa(int(*series.Id)),
			}, &episodeList); err != nil {
				//log.Printf("Failed to get episode list from Sonarr for '%s': %v", seriesTitle, err)
				continue
			}

			for _, ep := range episodeList {
				if ep.TvdbId != nil && *ep.TvdbId == episodeTvdbId {
					matchingEpisode = ep
					return
				}
			}
		}
	}

	return
}

var yearSuffixRegex = regexp.MustCompile(` \(\d{4}\)$`)

func cleanTitle(title string) string {
	return yearSuffixRegex.ReplaceAllString(title, "")
}

func ptr[T any](val T) *T {
	return &val
}

type sonarrAPIClient struct {
	url        string
	apiKey     string
	httpClient *http.Client
}

func newSonarrAPIClient(baseUrl *url.URL, apiKey string) *sonarrAPIClient {
	return &sonarrAPIClient{
		url:    baseUrl.JoinPath("api", "v3", "/").String(),
		apiKey: apiKey,
		httpClient: &http.Client{
			Transport: &http.Transport{
				Proxy:                 nil, // $HTTP_PROXY etc. ignored
				MaxIdleConns:          50,
				IdleConnTimeout:       http.DefaultTransport.(*http.Transport).IdleConnTimeout,
				TLSHandshakeTimeout:   http.DefaultTransport.(*http.Transport).TLSHandshakeTimeout,
				ExpectContinueTimeout: http.DefaultTransport.(*http.Transport).ExpectContinueTimeout,
				ResponseHeaderTimeout: 30 * time.Second,
				DialContext:           (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext,
				ForceAttemptHTTP2:     false,
			},
		},
	}
}

func (c *sonarrAPIClient) do(method string, endpoint string, queryParams map[string]string, reqBody any, respBody any) error {
	finalUrl := c.url + endpoint
	if queryParams != nil {
		u, err := url.Parse(finalUrl)
		if err != nil {
			return err
		}

		q := u.Query()
		for k, v := range queryParams {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()

		finalUrl = u.String()
	}

	var jsonBuf bytes.Buffer
	var pReqBody io.Reader
	if reqBody != nil {
		jsonEnc := json.NewEncoder(&jsonBuf)
		jsonEnc.SetEscapeHTML(false)
		if err := jsonEnc.Encode(reqBody); err != nil {
			return fmt.Errorf("failed to serialise request body to JSON for %s: %w", finalUrl, err)
		}
		pReqBody = &jsonBuf
	}

	req, err := http.NewRequest(method, finalUrl, pReqBody)
	if err != nil {
		return fmt.Errorf("failed to create %s request for %s: %w", method, finalUrl, err)
	}
	req.Header.Set("X-Api-Key", c.apiKey)
	if reqBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if respBody != nil {
		req.Header.Set("Accept", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute %s request for %s: %w", method, finalUrl, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= 300 {
		return fmt.Errorf("failed to %s %s: %s", method, finalUrl, resp.Status)
	}

	if respBody != nil {
		if err = json.NewDecoder(resp.Body).Decode(respBody); err != nil {
			return fmt.Errorf("failed to decode JSON response from %s: %w", finalUrl, err)
		}
	}

	return nil
}

func (c *sonarrAPIClient) get(endpoint string, queryParams map[string]string, respBody any) error {
	return c.do(http.MethodGet, endpoint, queryParams, nil, respBody)
}

func (c *sonarrAPIClient) put(endpoint string, queryParams map[string]string, reqBody any, respBody any) error {
	return c.do(http.MethodPut, endpoint, queryParams, reqBody, respBody)
}
