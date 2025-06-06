// Package sonarrt provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package sonarrt

import (
	"time"
)

const (
	X_Api_KeyScopes = "X_Api_Key.Scopes"
	ApikeyScopes    = "apikey.Scopes"
)

// Defines values for DownloadProtocol.
const (
	DownloadProtocolTorrent DownloadProtocol = "torrent"
	DownloadProtocolUnknown DownloadProtocol = "unknown"
	DownloadProtocolUsenet  DownloadProtocol = "usenet"
)

// Defines values for ImportListType.
const (
	ImportListTypeAdvanced ImportListType = "advanced"
	ImportListTypeOther    ImportListType = "other"
	ImportListTypePlex     ImportListType = "plex"
	ImportListTypeProgram  ImportListType = "program"
	ImportListTypeSimkl    ImportListType = "simkl"
	ImportListTypeTrakt    ImportListType = "trakt"
)

// Defines values for MediaCoverTypes.
const (
	MediaCoverTypesBanner     MediaCoverTypes = "banner"
	MediaCoverTypesClearlogo  MediaCoverTypes = "clearlogo"
	MediaCoverTypesFanart     MediaCoverTypes = "fanart"
	MediaCoverTypesHeadshot   MediaCoverTypes = "headshot"
	MediaCoverTypesPoster     MediaCoverTypes = "poster"
	MediaCoverTypesScreenshot MediaCoverTypes = "screenshot"
	MediaCoverTypesUnknown    MediaCoverTypes = "unknown"
)

// Defines values for MonitorTypes.
const (
	MonitorTypesAll               MonitorTypes = "all"
	MonitorTypesExisting          MonitorTypes = "existing"
	MonitorTypesFirstSeason       MonitorTypes = "firstSeason"
	MonitorTypesFuture            MonitorTypes = "future"
	MonitorTypesLastSeason        MonitorTypes = "lastSeason"
	MonitorTypesLatestSeason      MonitorTypes = "latestSeason"
	MonitorTypesMissing           MonitorTypes = "missing"
	MonitorTypesMonitorSpecials   MonitorTypes = "monitorSpecials"
	MonitorTypesNone              MonitorTypes = "none"
	MonitorTypesPilot             MonitorTypes = "pilot"
	MonitorTypesRecent            MonitorTypes = "recent"
	MonitorTypesSkip              MonitorTypes = "skip"
	MonitorTypesUnknown           MonitorTypes = "unknown"
	MonitorTypesUnmonitorSpecials MonitorTypes = "unmonitorSpecials"
)

// Defines values for NewItemMonitorTypes.
const (
	NewItemMonitorTypesAll  NewItemMonitorTypes = "all"
	NewItemMonitorTypesNone NewItemMonitorTypes = "none"
)

// Defines values for PrivacyLevel.
const (
	PrivacyLevelApiKey   PrivacyLevel = "apiKey"
	PrivacyLevelNormal   PrivacyLevel = "normal"
	PrivacyLevelPassword PrivacyLevel = "password"
	PrivacyLevelUserName PrivacyLevel = "userName"
)

// Defines values for ProviderMessageType.
const (
	ProviderMessageTypeError   ProviderMessageType = "error"
	ProviderMessageTypeInfo    ProviderMessageType = "info"
	ProviderMessageTypeWarning ProviderMessageType = "warning"
)

// Defines values for QualitySource.
const (
	QualitySourceBluray        QualitySource = "bluray"
	QualitySourceBlurayRaw     QualitySource = "blurayRaw"
	QualitySourceDvd           QualitySource = "dvd"
	QualitySourceTelevision    QualitySource = "television"
	QualitySourceTelevisionRaw QualitySource = "televisionRaw"
	QualitySourceUnknown       QualitySource = "unknown"
	QualitySourceWeb           QualitySource = "web"
	QualitySourceWebRip        QualitySource = "webRip"
)

// Defines values for SeriesStatusType.
const (
	SeriesStatusTypeContinuing SeriesStatusType = "continuing"
	SeriesStatusTypeDeleted    SeriesStatusType = "deleted"
	SeriesStatusTypeEnded      SeriesStatusType = "ended"
	SeriesStatusTypeUpcoming   SeriesStatusType = "upcoming"
)

// Defines values for SeriesTypes.
const (
	SeriesTypesAnime    SeriesTypes = "anime"
	SeriesTypesDaily    SeriesTypes = "daily"
	SeriesTypesStandard SeriesTypes = "standard"
)

// Defines values for SortDirection.
const (
	SortDirectionAscending  SortDirection = "ascending"
	SortDirectionDefault    SortDirection = "default"
	SortDirectionDescending SortDirection = "descending"
)

// AddSeriesOptions defines model for AddSeriesOptions.
type AddSeriesOptions struct {
	IgnoreEpisodesWithFiles      *bool         `json:"ignoreEpisodesWithFiles,omitempty"`
	IgnoreEpisodesWithoutFiles   *bool         `json:"ignoreEpisodesWithoutFiles,omitempty"`
	Monitor                      *MonitorTypes `json:"monitor,omitempty"`
	SearchForCutoffUnmetEpisodes *bool         `json:"searchForCutoffUnmetEpisodes,omitempty"`
	SearchForMissingEpisodes     *bool         `json:"searchForMissingEpisodes,omitempty"`
}

// AlternateTitleResource defines model for AlternateTitleResource.
type AlternateTitleResource struct {
	Comment           *string `json:"comment"`
	SceneOrigin       *string `json:"sceneOrigin"`
	SceneSeasonNumber *int32  `json:"sceneSeasonNumber"`
	SeasonNumber      *int32  `json:"seasonNumber"`
	Title             *string `json:"title"`
}

// BlocklistBulkResource defines model for BlocklistBulkResource.
type BlocklistBulkResource struct {
	Ids *[]int32 `json:"ids"`
}

// BlocklistResource defines model for BlocklistResource.
type BlocklistResource struct {
	CustomFormats *[]CustomFormatResource `json:"customFormats"`
	Date          *time.Time              `json:"date,omitempty"`
	EpisodeIds    *[]int32                `json:"episodeIds"`
	Id            *int32                  `json:"id,omitempty"`
	Indexer       *string                 `json:"indexer"`
	Languages     *[]Language             `json:"languages"`
	Message       *string                 `json:"message"`
	Protocol      *DownloadProtocol       `json:"protocol,omitempty"`
	Quality       *QualityModel           `json:"quality,omitempty"`
	Series        *SeriesResource         `json:"series,omitempty"`
	SeriesId      *int32                  `json:"seriesId,omitempty"`
	SourceTitle   *string                 `json:"sourceTitle"`
}

// BlocklistResourcePagingResource defines model for BlocklistResourcePagingResource.
type BlocklistResourcePagingResource struct {
	Page          *int32               `json:"page,omitempty"`
	PageSize      *int32               `json:"pageSize,omitempty"`
	Records       *[]BlocklistResource `json:"records"`
	SortDirection *SortDirection       `json:"sortDirection,omitempty"`
	SortKey       *string              `json:"sortKey"`
	TotalRecords  *int32               `json:"totalRecords,omitempty"`
}

// CustomFormatResource defines model for CustomFormatResource.
type CustomFormatResource struct {
	Id                              *int32                             `json:"id,omitempty"`
	IncludeCustomFormatWhenRenaming *bool                              `json:"includeCustomFormatWhenRenaming"`
	Name                            *string                            `json:"name"`
	Specifications                  *[]CustomFormatSpecificationSchema `json:"specifications"`
}

// CustomFormatSpecificationSchema defines model for CustomFormatSpecificationSchema.
type CustomFormatSpecificationSchema struct {
	Fields             *[]Field                           `json:"fields"`
	Id                 *int32                             `json:"id,omitempty"`
	Implementation     *string                            `json:"implementation"`
	ImplementationName *string                            `json:"implementationName"`
	InfoLink           *string                            `json:"infoLink"`
	Name               *string                            `json:"name"`
	Negate             *bool                              `json:"negate,omitempty"`
	Presets            *[]CustomFormatSpecificationSchema `json:"presets"`
	Required           *bool                              `json:"required,omitempty"`
}

// DownloadClientResource defines model for DownloadClientResource.
type DownloadClientResource struct {
	ConfigContract           *string                   `json:"configContract"`
	Enable                   *bool                     `json:"enable,omitempty"`
	Fields                   *[]Field                  `json:"fields"`
	Id                       *int32                    `json:"id,omitempty"`
	Implementation           *string                   `json:"implementation"`
	ImplementationName       *string                   `json:"implementationName"`
	InfoLink                 *string                   `json:"infoLink"`
	Message                  *ProviderMessage          `json:"message,omitempty"`
	Name                     *string                   `json:"name"`
	Presets                  *[]DownloadClientResource `json:"presets"`
	Priority                 *int32                    `json:"priority,omitempty"`
	Protocol                 *DownloadProtocol         `json:"protocol,omitempty"`
	RemoveCompletedDownloads *bool                     `json:"removeCompletedDownloads,omitempty"`
	RemoveFailedDownloads    *bool                     `json:"removeFailedDownloads,omitempty"`
	Tags                     *[]int32                  `json:"tags"`
}

// DownloadProtocol defines model for DownloadProtocol.
type DownloadProtocol string

// EpisodeResource defines model for EpisodeResource.
type EpisodeResource struct {
	EpisodeNumber *int32          `json:"episodeNumber,omitempty"`
	Id            *int32          `json:"id,omitempty"`
	Monitored     *bool           `json:"monitored,omitempty"`
	SeasonNumber  *int32          `json:"seasonNumber,omitempty"`
	Series        *SeriesResource `json:"series,omitempty"`
	TvdbId        *int32          `json:"tvdbId,omitempty"`
}

// EpisodesMonitoredResource defines model for EpisodesMonitoredResource.
type EpisodesMonitoredResource struct {
	EpisodeIds *[]int32 `json:"episodeIds"`
	Monitored  *bool    `json:"monitored,omitempty"`
}

// Field defines model for Field.
type Field struct {
	Advanced                    *bool           `json:"advanced,omitempty"`
	HelpLink                    *string         `json:"helpLink"`
	HelpText                    *string         `json:"helpText"`
	HelpTextWarning             *string         `json:"helpTextWarning"`
	Hidden                      *string         `json:"hidden"`
	IsFloat                     *bool           `json:"isFloat,omitempty"`
	Label                       *string         `json:"label"`
	Name                        *string         `json:"name"`
	Order                       *int32          `json:"order,omitempty"`
	Placeholder                 *string         `json:"placeholder"`
	Privacy                     *PrivacyLevel   `json:"privacy,omitempty"`
	Section                     *string         `json:"section"`
	SelectOptions               *[]SelectOption `json:"selectOptions"`
	SelectOptionsProviderAction *string         `json:"selectOptionsProviderAction"`
	Type                        *string         `json:"type"`
	Unit                        *string         `json:"unit"`
	Value                       *interface{}    `json:"value"`
}

// ImportListResource defines model for ImportListResource.
type ImportListResource struct {
	ConfigContract           *string               `json:"configContract"`
	EnableAutomaticAdd       *bool                 `json:"enableAutomaticAdd,omitempty"`
	Fields                   *[]Field              `json:"fields"`
	Id                       *int32                `json:"id,omitempty"`
	Implementation           *string               `json:"implementation"`
	ImplementationName       *string               `json:"implementationName"`
	InfoLink                 *string               `json:"infoLink"`
	ListOrder                *int32                `json:"listOrder,omitempty"`
	ListType                 *ImportListType       `json:"listType,omitempty"`
	Message                  *ProviderMessage      `json:"message,omitempty"`
	MinRefreshInterval       *string               `json:"minRefreshInterval,omitempty"`
	MonitorNewItems          *NewItemMonitorTypes  `json:"monitorNewItems,omitempty"`
	Name                     *string               `json:"name"`
	Presets                  *[]ImportListResource `json:"presets"`
	QualityProfileId         *int32                `json:"qualityProfileId,omitempty"`
	RootFolderPath           *string               `json:"rootFolderPath"`
	SearchForMissingEpisodes *bool                 `json:"searchForMissingEpisodes,omitempty"`
	SeasonFolder             *bool                 `json:"seasonFolder,omitempty"`
	SeriesType               *SeriesTypes          `json:"seriesType,omitempty"`
	ShouldMonitor            *MonitorTypes         `json:"shouldMonitor,omitempty"`
	Tags                     *[]int32              `json:"tags"`
}

// ImportListType defines model for ImportListType.
type ImportListType string

// IndexerResource defines model for IndexerResource.
type IndexerResource struct {
	ConfigContract                      *string            `json:"configContract"`
	DownloadClientId                    *int32             `json:"downloadClientId,omitempty"`
	EnableAutomaticSearch               *bool              `json:"enableAutomaticSearch,omitempty"`
	EnableInteractiveSearch             *bool              `json:"enableInteractiveSearch,omitempty"`
	EnableRss                           *bool              `json:"enableRss,omitempty"`
	Fields                              *[]Field           `json:"fields"`
	Id                                  *int32             `json:"id,omitempty"`
	Implementation                      *string            `json:"implementation"`
	ImplementationName                  *string            `json:"implementationName"`
	InfoLink                            *string            `json:"infoLink"`
	Message                             *ProviderMessage   `json:"message,omitempty"`
	Name                                *string            `json:"name"`
	Presets                             *[]IndexerResource `json:"presets"`
	Priority                            *int32             `json:"priority,omitempty"`
	Protocol                            *DownloadProtocol  `json:"protocol,omitempty"`
	SeasonSearchMaximumSingleEpisodeAge *int32             `json:"seasonSearchMaximumSingleEpisodeAge,omitempty"`
	SupportsRss                         *bool              `json:"supportsRss,omitempty"`
	SupportsSearch                      *bool              `json:"supportsSearch,omitempty"`
	Tags                                *[]int32           `json:"tags"`
}

// Language defines model for Language.
type Language struct {
	Id   *int32  `json:"id,omitempty"`
	Name *string `json:"name"`
}

// MediaCover defines model for MediaCover.
type MediaCover struct {
	CoverType *MediaCoverTypes `json:"coverType,omitempty"`
	RemoteUrl *string          `json:"remoteUrl"`
	Url       *string          `json:"url"`
}

// MediaCoverTypes defines model for MediaCoverTypes.
type MediaCoverTypes string

// MetadataResource defines model for MetadataResource.
type MetadataResource struct {
	ConfigContract     *string             `json:"configContract"`
	Enable             *bool               `json:"enable,omitempty"`
	Fields             *[]Field            `json:"fields"`
	Id                 *int32              `json:"id,omitempty"`
	Implementation     *string             `json:"implementation"`
	ImplementationName *string             `json:"implementationName"`
	InfoLink           *string             `json:"infoLink"`
	Message            *ProviderMessage    `json:"message,omitempty"`
	Name               *string             `json:"name"`
	Presets            *[]MetadataResource `json:"presets"`
	Tags               *[]int32            `json:"tags"`
}

// MonitorTypes defines model for MonitorTypes.
type MonitorTypes string

// NewItemMonitorTypes defines model for NewItemMonitorTypes.
type NewItemMonitorTypes string

// NotificationResource defines model for NotificationResource.
type NotificationResource struct {
	ConfigContract                        *string                 `json:"configContract"`
	Fields                                *[]Field                `json:"fields"`
	Id                                    *int32                  `json:"id,omitempty"`
	Implementation                        *string                 `json:"implementation"`
	ImplementationName                    *string                 `json:"implementationName"`
	IncludeHealthWarnings                 *bool                   `json:"includeHealthWarnings,omitempty"`
	InfoLink                              *string                 `json:"infoLink"`
	Link                                  *string                 `json:"link"`
	Message                               *ProviderMessage        `json:"message,omitempty"`
	Name                                  *string                 `json:"name"`
	OnApplicationUpdate                   *bool                   `json:"onApplicationUpdate,omitempty"`
	OnDownload                            *bool                   `json:"onDownload,omitempty"`
	OnEpisodeFileDelete                   *bool                   `json:"onEpisodeFileDelete,omitempty"`
	OnEpisodeFileDeleteForUpgrade         *bool                   `json:"onEpisodeFileDeleteForUpgrade,omitempty"`
	OnGrab                                *bool                   `json:"onGrab,omitempty"`
	OnHealthIssue                         *bool                   `json:"onHealthIssue,omitempty"`
	OnHealthRestored                      *bool                   `json:"onHealthRestored,omitempty"`
	OnImportComplete                      *bool                   `json:"onImportComplete,omitempty"`
	OnManualInteractionRequired           *bool                   `json:"onManualInteractionRequired,omitempty"`
	OnRename                              *bool                   `json:"onRename,omitempty"`
	OnSeriesAdd                           *bool                   `json:"onSeriesAdd,omitempty"`
	OnSeriesDelete                        *bool                   `json:"onSeriesDelete,omitempty"`
	OnUpgrade                             *bool                   `json:"onUpgrade,omitempty"`
	Presets                               *[]NotificationResource `json:"presets"`
	SupportsOnApplicationUpdate           *bool                   `json:"supportsOnApplicationUpdate,omitempty"`
	SupportsOnDownload                    *bool                   `json:"supportsOnDownload,omitempty"`
	SupportsOnEpisodeFileDelete           *bool                   `json:"supportsOnEpisodeFileDelete,omitempty"`
	SupportsOnEpisodeFileDeleteForUpgrade *bool                   `json:"supportsOnEpisodeFileDeleteForUpgrade,omitempty"`
	SupportsOnGrab                        *bool                   `json:"supportsOnGrab,omitempty"`
	SupportsOnHealthIssue                 *bool                   `json:"supportsOnHealthIssue,omitempty"`
	SupportsOnHealthRestored              *bool                   `json:"supportsOnHealthRestored,omitempty"`
	SupportsOnImportComplete              *bool                   `json:"supportsOnImportComplete,omitempty"`
	SupportsOnManualInteractionRequired   *bool                   `json:"supportsOnManualInteractionRequired,omitempty"`
	SupportsOnRename                      *bool                   `json:"supportsOnRename,omitempty"`
	SupportsOnSeriesAdd                   *bool                   `json:"supportsOnSeriesAdd,omitempty"`
	SupportsOnSeriesDelete                *bool                   `json:"supportsOnSeriesDelete,omitempty"`
	SupportsOnUpgrade                     *bool                   `json:"supportsOnUpgrade,omitempty"`
	Tags                                  *[]int32                `json:"tags"`
	TestCommand                           *string                 `json:"testCommand"`
}

// PrivacyLevel defines model for PrivacyLevel.
type PrivacyLevel string

// ProviderMessage defines model for ProviderMessage.
type ProviderMessage struct {
	Message *string              `json:"message"`
	Type    *ProviderMessageType `json:"type,omitempty"`
}

// ProviderMessageType defines model for ProviderMessageType.
type ProviderMessageType string

// Quality defines model for Quality.
type Quality struct {
	Id         *int32         `json:"id,omitempty"`
	Name       *string        `json:"name"`
	Resolution *int32         `json:"resolution,omitempty"`
	Source     *QualitySource `json:"source,omitempty"`
}

// QualityModel defines model for QualityModel.
type QualityModel struct {
	Quality  *Quality  `json:"quality,omitempty"`
	Revision *Revision `json:"revision,omitempty"`
}

// QualityProfileQualityItemResource defines model for QualityProfileQualityItemResource.
type QualityProfileQualityItemResource struct {
	Allowed *bool                                `json:"allowed,omitempty"`
	Id      *int32                               `json:"id,omitempty"`
	Items   *[]QualityProfileQualityItemResource `json:"items"`
	Name    *string                              `json:"name"`
	Quality *Quality                             `json:"quality,omitempty"`
}

// QualitySource defines model for QualitySource.
type QualitySource string

// Ratings defines model for Ratings.
type Ratings struct {
	Value *float64 `json:"value,omitempty"`
	Votes *int32   `json:"votes,omitempty"`
}

// Revision defines model for Revision.
type Revision struct {
	IsRepack *bool  `json:"isRepack,omitempty"`
	Real     *int32 `json:"real,omitempty"`
	Version  *int32 `json:"version,omitempty"`
}

// RootFolderResource defines model for RootFolderResource.
type RootFolderResource struct {
	Path *string `json:"path"`
}

// SeasonResource defines model for SeasonResource.
type SeasonResource struct {
	Images       *[]MediaCover             `json:"images"`
	Monitored    *bool                     `json:"monitored,omitempty"`
	SeasonNumber *int32                    `json:"seasonNumber,omitempty"`
	Statistics   *SeasonStatisticsResource `json:"statistics,omitempty"`
}

// SeasonStatisticsResource defines model for SeasonStatisticsResource.
type SeasonStatisticsResource struct {
	EpisodeCount      *int32     `json:"episodeCount,omitempty"`
	EpisodeFileCount  *int32     `json:"episodeFileCount,omitempty"`
	NextAiring        *time.Time `json:"nextAiring"`
	PercentOfEpisodes *float64   `json:"percentOfEpisodes,omitempty"`
	PreviousAiring    *time.Time `json:"previousAiring"`
	ReleaseGroups     *[]string  `json:"releaseGroups"`
	SizeOnDisk        *int64     `json:"sizeOnDisk,omitempty"`
	TotalEpisodeCount *int32     `json:"totalEpisodeCount,omitempty"`
}

// SelectOption defines model for SelectOption.
type SelectOption struct {
	Hint  *string `json:"hint"`
	Name  *string `json:"name"`
	Order *int32  `json:"order,omitempty"`
	Value *int32  `json:"value,omitempty"`
}

// SeriesResource defines model for SeriesResource.
type SeriesResource struct {
	AddOptions      *AddSeriesOptions         `json:"addOptions,omitempty"`
	Added           *time.Time                `json:"added,omitempty"`
	AirTime         *string                   `json:"airTime"`
	AlternateTitles *[]AlternateTitleResource `json:"alternateTitles"`
	Certification   *string                   `json:"certification"`
	CleanTitle      *string                   `json:"cleanTitle"`
	Ended           *bool                     `json:"ended,omitempty"`
	EpisodesChanged *bool                     `json:"episodesChanged"`
	FirstAired      *time.Time                `json:"firstAired"`
	Folder          *string                   `json:"folder"`
	Genres          *[]string                 `json:"genres"`
	Id              *int32                    `json:"id,omitempty"`
	Images          *[]MediaCover             `json:"images"`
	ImdbId          *string                   `json:"imdbId"`
	// Deprecated:
	LanguageProfileId *int32                    `json:"languageProfileId,omitempty"`
	LastAired         *time.Time                `json:"lastAired"`
	MonitorNewItems   *NewItemMonitorTypes      `json:"monitorNewItems,omitempty"`
	Monitored         *bool                     `json:"monitored,omitempty"`
	Network           *string                   `json:"network"`
	NextAiring        *time.Time                `json:"nextAiring"`
	OriginalLanguage  *Language                 `json:"originalLanguage,omitempty"`
	Overview          *string                   `json:"overview"`
	Path              *string                   `json:"path"`
	PreviousAiring    *time.Time                `json:"previousAiring"`
	ProfileName       *string                   `json:"profileName"`
	QualityProfileId  *int32                    `json:"qualityProfileId,omitempty"`
	Ratings           *Ratings                  `json:"ratings,omitempty"`
	RemotePoster      *string                   `json:"remotePoster"`
	RootFolderPath    *string                   `json:"rootFolderPath"`
	Runtime           *int32                    `json:"runtime,omitempty"`
	SeasonFolder      *bool                     `json:"seasonFolder,omitempty"`
	Seasons           *[]SeasonResource         `json:"seasons"`
	SeriesType        *SeriesTypes              `json:"seriesType,omitempty"`
	SortTitle         *string                   `json:"sortTitle"`
	Statistics        *SeriesStatisticsResource `json:"statistics,omitempty"`
	Status            *SeriesStatusType         `json:"status,omitempty"`
	Tags              *[]int32                  `json:"tags"`
	Title             *string                   `json:"title"`
	TitleSlug         *string                   `json:"titleSlug"`
	TmdbId            *int32                    `json:"tmdbId,omitempty"`
	TvMazeId          *int32                    `json:"tvMazeId,omitempty"`
	TvRageId          *int32                    `json:"tvRageId,omitempty"`
	TvdbId            *int32                    `json:"tvdbId,omitempty"`
	UseSceneNumbering *bool                     `json:"useSceneNumbering,omitempty"`
	Year              *int32                    `json:"year,omitempty"`
}

// SeriesStatisticsResource defines model for SeriesStatisticsResource.
type SeriesStatisticsResource struct {
	EpisodeCount      *int32    `json:"episodeCount,omitempty"`
	EpisodeFileCount  *int32    `json:"episodeFileCount,omitempty"`
	PercentOfEpisodes *float64  `json:"percentOfEpisodes,omitempty"`
	ReleaseGroups     *[]string `json:"releaseGroups"`
	SeasonCount       *int32    `json:"seasonCount,omitempty"`
	SizeOnDisk        *int64    `json:"sizeOnDisk,omitempty"`
	TotalEpisodeCount *int32    `json:"totalEpisodeCount,omitempty"`
}

// SeriesStatusType defines model for SeriesStatusType.
type SeriesStatusType string

// SeriesTypes defines model for SeriesTypes.
type SeriesTypes string

// SortDirection defines model for SortDirection.
type SortDirection string

// GetApiV3BlocklistParams defines parameters for GetApiV3Blocklist.
type GetApiV3BlocklistParams struct {
	Page          *int32              `form:"page,omitempty" json:"page,omitempty"`
	PageSize      *int32              `form:"pageSize,omitempty" json:"pageSize,omitempty"`
	SortKey       *string             `form:"sortKey,omitempty" json:"sortKey,omitempty"`
	SortDirection *SortDirection      `form:"sortDirection,omitempty" json:"sortDirection,omitempty"`
	SeriesIds     *[]int32            `form:"seriesIds,omitempty" json:"seriesIds,omitempty"`
	Protocols     *[]DownloadProtocol `form:"protocols,omitempty" json:"protocols,omitempty"`
}

// GetApiV3EpisodeParams defines parameters for GetApiV3Episode.
type GetApiV3EpisodeParams struct {
	SeriesId           *int32   `form:"seriesId,omitempty" json:"seriesId,omitempty"`
	SeasonNumber       *int32   `form:"seasonNumber,omitempty" json:"seasonNumber,omitempty"`
	EpisodeIds         *[]int32 `form:"episodeIds,omitempty" json:"episodeIds,omitempty"`
	EpisodeFileId      *int32   `form:"episodeFileId,omitempty" json:"episodeFileId,omitempty"`
	IncludeSeries      *bool    `form:"includeSeries,omitempty" json:"includeSeries,omitempty"`
	IncludeEpisodeFile *bool    `form:"includeEpisodeFile,omitempty" json:"includeEpisodeFile,omitempty"`
	IncludeImages      *bool    `form:"includeImages,omitempty" json:"includeImages,omitempty"`
}

// PutApiV3EpisodeMonitorParams defines parameters for PutApiV3EpisodeMonitor.
type PutApiV3EpisodeMonitorParams struct {
	IncludeImages *bool `form:"includeImages,omitempty" json:"includeImages,omitempty"`
}

// DeleteApiV3BlocklistBulkApplicationWildcardPlusJSONRequestBody defines body for DeleteApiV3BlocklistBulk for application/*+json ContentType.
type DeleteApiV3BlocklistBulkApplicationWildcardPlusJSONRequestBody = BlocklistBulkResource

// DeleteApiV3BlocklistBulkJSONRequestBody defines body for DeleteApiV3BlocklistBulk for application/json ContentType.
type DeleteApiV3BlocklistBulkJSONRequestBody = BlocklistBulkResource

// PutApiV3EpisodeMonitorJSONRequestBody defines body for PutApiV3EpisodeMonitor for application/json ContentType.
type PutApiV3EpisodeMonitorJSONRequestBody = EpisodesMonitoredResource

// PutApiV3EpisodeIdJSONRequestBody defines body for PutApiV3EpisodeId for application/json ContentType.
type PutApiV3EpisodeIdJSONRequestBody = EpisodeResource

// PostApiV3RootfolderJSONRequestBody defines body for PostApiV3Rootfolder for application/json ContentType.
type PostApiV3RootfolderJSONRequestBody = RootFolderResource
