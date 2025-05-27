package sonarrt

//go:generate go tool -modfile=../../go.tool.mod oapi-codegen -config cfg.yaml https://raw.githubusercontent.com/Sonarr/Sonarr/develop/src/Sonarr.Api.V3/openapi.json
