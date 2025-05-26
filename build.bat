@echo off
setlocal
pushd "%~dp0"

REM First off, this is indeed a batch file... And it's rather opinionated for use on my system. Using the usual go build commands instead, invoked standalone, should be fine with safe defaults
set GOTELEMETRY=off
set GOAMD64=v3
::go generate ./... || exit /b %ERRORLEVEL%
go build -trimpath -gcflags="all=-C -dwarf=false" -ldflags="-s -w -buildid=" -buildvcs=false -tags=mapof_opt_cachelinesize_64,mapof_opt_enablepadding,mapof_opt_atomiclevel_1

popd
