APP := gotp
GO_LDFLAGS="-s -extldflags=-static"

_: rmbuild linux darwin windows
	tar cfv build/gotp-linux_v0.0.1.tar build/$(APP)-linux
	tar cfv build/gotp-darwin_v0.0.1.tar build/$(APP)-darwin
	zip build/gotp-windows_v0.0.1.zip build/$(APP).exe

rmbuild:
	rm -rf build

clean:
	@echo "Building ..."
	go clean
	go get	

linux: clean
	@GOOS=linux CGO_ENABLED=0 go build -ldflags=$(GO_LDFLAGS) -o build/$(APP)-linux
	@echo "finished"

darwin: clean
	@GOOS=darwin go build -ldflags=$(GO_LDFLAGS) -o build/$(APP)-darwin
	@echo "finished"

windows: clean
	@GOOS=windows go build -ldflags=$(GO_LDFLAGS) -o build/$(APP).exe
	@echo "finished"

