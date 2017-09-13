all: main

main:
	@mkdir -p dist/bin dist/conf dist/log
	@echo "Release version"
	@${GOROOT}/bin/go build -v
	@mv macedon dist/bin/macedon
	@cp conf/macedon.conf dist/conf/
	@echo "Build done: binary in dist dir"
debug:
	@mkdir -p dist/bin dist/conf dist/log
	@echo "Debug version"
	@${GOROOT}/bin/go build -o dist/bin/macedon -ldflags '-s -w' main.go
	@cp conf/macedon.conf dist/conf/
	@echo "Build done: binary in dist dir"

clean:
	@rm -rf dist

.PHONY: all main clean debug
