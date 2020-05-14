#! /bin/sh

VERSION="$1"

mkdir build

go clean
mkdir build/win32
make build SYSTEM="GOOS=windows GOARCH=386"
ln ddns.exe build/win32/
zip -9 -r ddns-win32-${VERSION:-dev}.zip build/win32

go clean
mkdir build/win64
make build SYSTEM="GOOS=windows GOARCH=amd64"
ln ddns.exe build/win64/
zip -9 -r ddns-win64-${VERSION:-dev}.zip build/win64