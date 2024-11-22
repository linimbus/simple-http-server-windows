rsrc -manifest exe.manifest -ico static/main.ico
rice embed-go
go build -buildvcs=false -ldflags="-H windowsgui -w -s" -o SimpleHttpFileServer.exe