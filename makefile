macos_ENV=CGO_ENABLED=0 GOOS=darwin GOARCH=amd64
windows_ENV=CGO_ENABLED=0 GOOS=windows GOARCH=amd64
linux_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64

win: PLAT=windows
mac: PLAT=macos
linux: PLAT=linux

win: all
mac: all
linux: all

all: ; $($(PLAT)_ENV) go build
