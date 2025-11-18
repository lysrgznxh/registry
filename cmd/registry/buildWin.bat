@ECHO OFF
set GOOS=windows
set CGO_ENABLED=1
go build -buildmode exe -o registry.exe -ldflags="-s -w"
