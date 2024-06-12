set GOOS=linux& set GOARCH=arm& go build -v .
ren engine IdleReader_v0.2.0_mac_linux

set GOOS=windows& set GOARCH=amd64& go build -v .
ren engine.exe IdleReader_v0.2.0_windows.exe