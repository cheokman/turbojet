export VERSION=0.1.0
export RELEASE_PATH="releases/turbojet-${VERSION}"
export DRIVER_PATH="driver"
export RESOURCES_PATH="resources"

all: build

deps: 
	-go get -t -d github.com/tebeka/selenium
	dep ensure

build: deps
  go build -ldflag "-X 'turbojet/cli.Version=${VERSION}'" -o out/turbojet/bin/tj main/main.go 

release: mac linux windows

mac: clean make_build_dir make_release_dir release_mac

linux: clean make_build_dir make_release_dir release_linux

windows: clean make_build_dir make_release_dir release_windows

make_build_dir:
	mkdir -p out/turbojet/driver

make_release_dir:
	mkdir -p ${RELEASE_PATH}

clean:
	rm -rf out/*

release_mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'turbojet/cli.Version=${VERSION}'" -o out/turbojet/bin/tj main/main.go
	cp ${DRIVER_PATH}/chromedriver-mac out/turbojet/driver/chromedriver
	cp ${DRIVER_PATH}/selenium-server-standalone-3.141.59.jar out/turbojet/driver/selenium-server-standalone-3.141.59.jar
	cp -R ${RESOURCES_PATH}/content out/turbojet/
	tar zcvf ${RELEASE_PATH}/turbojet-darwin-amd64-${VERSION}.tar.gz -C out turbojet

release_linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-X 'turbojet/cli.Version=${VERSION}'" -o out/turbojet/bin/tj main/main.go
	cp ${DRIVER_PATH}/chromedriver-linux out/turbojet/driver/chromedriver
	cp ${DRIVER_PATH}/google-chrome-stable-73.0.3683.75-1.x86_64.rpm out/turbojet/driver/google-chrome-stable-73.0.3683.75-1.x86_64.rpm
	cp ${DRIVER_PATH}/selenium-server-standalone-3.141.59.jar out/turbojet/driver/selenium-server-standalone-3.141.59.jar
	cp -R ${RESOURCES_PATH}/content out/turbojet/
	tar zcvf ${RELEASE_PATH}/turbojet-linux-amd64-${VERSION}.tar.gz -C out turbojet

release_windows:
	GOOS=windows GOARCH=amd64 go build -ldflags "-X 'turbojet/cli.Version=${VERSION}'" -o out/turbojet/bin/tj.exe main/main.go
	cp ${DRIVER_PATH}/chromedriver-win.exe out/turbojet/driver/chromedriver.exe
	cp -R ${RESOURCES_PATH}/content out/turbojet/
	cp ${DRIVER_PATH}/selenium-server-standalone-3.141.59.jar out/turbojet/driver/selenium-server-standalone-3.141.59.jar
	zip -r ${RELEASE_PATH}/turbojet-windows-amd64.exe.zip out/turbojet