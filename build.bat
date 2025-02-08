rsrc -manifest exe.manifest -ico main.ico
go-bindata -o icon_files.go main.ico
go build -ldflags="-H windowsgui -w -s" -o simple-PNG-ICO-windows-x64.exe
