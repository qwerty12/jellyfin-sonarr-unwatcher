@echo off
setlocal
pushd "%~dp0"

REM This is indeed a batch file. And it's rather opinionated being intended for use on my system. Run the usual go build commands instead, standalone, for safe defaults
set GOTELEMETRY=off
set GOAMD64=v3
set GOEXPERIMENT=jsonv2,greenteagc
::go generate ./... || exit /b %ERRORLEVEL%
go build -trimpath -gcflags="all=-C -dwarf=false" -ldflags="-s -w -buildid=" -buildvcs=false -tags=mapof_opt_cachelinesize_64,mapof_opt_enablepadding,mapof_opt_atomiclevel_2

popd
