default: build
PKG_CONFIG := $(shell pwd)

clean :
	rm -fr dist
init:
	mkdir dist -p
	go mod tidy
build:
	make clean
	make init
	gomobile bind -v -o dist/lib_onvif.aar -target=android .
build_ios:
	make clean
	make init
	gomobile bind -v -o dist/lib_onvif.aar -target=ios .


