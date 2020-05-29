mac_ENV=CGO_ENABLED=0 GOOS=darwin GOARCH=amd64
win_ENV=CGO_ENABLED=0 GOOS=windows GOARCH=amd64
lin_ENV=CGO_ENABLED=0 GOOS=linux GOARCH=amd64

windows: PLAT=win
macos: PLAT=mac
linux: PLAT=lin

windows: all
macos: all
linux: all

all: ;  $($(PLAT)_ENV) go build
