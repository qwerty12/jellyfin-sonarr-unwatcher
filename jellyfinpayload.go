// qwerty12: This file was initially generated from JSON Schema using quicktype, and then referenced against the
// source of shemanaev's Webhook plugin.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    jellyfinPayload, err := UnmarshalJellyfinPayload(bytes)
//    bytes, err = jellyfinPayload.Marshal()

package main

import (
	"encoding/json"
	"io"
	"jellyfin-sonarr-unwatcher/internal/jellygen"
)

func DecodeJellyfinPayload(data io.Reader) (r JellyfinPayload, err error) {
	err = json.NewDecoder(data).Decode(&r)
	return
}

/*func UnmarshalJellyfinPayload(data []byte) (JellyfinPayload, error) {
	var r JellyfinPayload
	err := json.Unmarshal(data, &r)
	return r, err
}*/

func (r *JellyfinPayload) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// DefaultFormatPayload: https://github.com/shemanaev/jellyfin-plugin-webhooks/blob/master/Jellyfin.Webhooks/Formats/DefaultFormat.cs#L50
type JellyfinPayload struct {
	//Event          HookEvent             `json:"Event"`
	Item *jellygen.BaseItemDto `json:"Item,omitempty,omitzero"`
	//User *jellygen.UserDto     `json:"User,omitempty,omitzero"`
	//Session        *SessionInfoDto       `json:"Session,omitempty,omitzero"`
	//Server         ServerInfoDto         `json:"Server"`
	//AdditionalData *any                  `json:"AdditionalData,omitempty,omitzero"`
	Series *jellygen.BaseItemDto `json:"Series,omitempty,omitzero"`
}

/*
type HookEvent string

const (
	EventPlay                     HookEvent = "Play"
	EventPause                    HookEvent = "Pause"
	EventResume                   HookEvent = "Resume"
	EventStop                     HookEvent = "Stop"
	EventScrobble                 HookEvent = "Scrobble"
	EventProgress                 HookEvent = "Progress"
	EventMarkPlayed               HookEvent = "MarkPlayed"
	EventMarkUnplayed             HookEvent = "MarkUnplayed"
	EventRate                     HookEvent = "Rate"
	EventItemAdded                HookEvent = "ItemAdded"
	EventItemRemoved              HookEvent = "ItemRemoved"
	EventItemUpdated              HookEvent = "ItemUpdated"
	EventAuthenticationSucceeded  HookEvent = "AuthenticationSucceeded"
	EventAuthenticationFailed     HookEvent = "AuthenticationFailed"
	EventSessionStarted           HookEvent = "SessionStarted"
	EventSessionEnded             HookEvent = "SessionEnded"
	EventSubtitleDownloadFailure  HookEvent = "SubtitleDownloadFailure"
	EventHasPendingRestartChanged HookEvent = "HasPendingRestartChanged"
)

type ServerInfoDto struct {
	ID      *string `json:"Id,omitempty,omitzero"`
	Name    string  `json:"Name"`
	Version string  `json:"Version"`
}

type SessionInfoDto struct {
	ID                 *string                   `json:"Id,omitempty,omitzero"`
	Client             *string                   `json:"Client,omitempty,omitzero"`
	DeviceID           *string                   `json:"DeviceId,omitempty,omitzero"`
	DeviceName         *string                   `json:"DeviceName,omitempty,omitzero"`
	RemoteEndPoint     *string                   `json:"RemoteEndPoint,omitempty,omitzero"`
	ApplicationVersion *string                   `json:"ApplicationVersion,omitempty,omitzero"`
	PlayState          *jellygen.PlayerStateInfo `json:"PlayState,omitempty,omitzero"`
}
*/
