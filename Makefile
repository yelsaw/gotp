APP := gotp
GO_LDFLAGS="-s -extldflags=-static"

_: linux darwin windows
	tar cfv gotp-linux_v0.0.1.tar build/linux/*
	tar cfv gotp-darwin_v0.0.1.tar build/darwin/*
	zip gotp-windows_v0.0.1.zip build/windows/*

clean:
	rm -rf build

linux:
	@echo "Building ..."
	go clean
	go get
	@GOOS=linux CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o build/linux/$(APP)-linux
	@echo "finished"

darwin:
	@echo "Building ..."
	go clean
	go get
	@GOOS=darwin go build -ldflags=$(GO_LDFLAGS) -o build/darwin/$(APP)-darwin
	@echo "finished"

windows:
	@echo "Building ..."
	go clean
	go get
	@GOOS=windows go build -ldflags=$(GO_LDFLAGS) -o build/windows/$(APP).exe
	@echo "finished"

