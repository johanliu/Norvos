all: main

main:
	@go build --buildmode=plugin -o /opt/essos/dns.so dns.go

clean:
	@rm -rf dist

.PHONY: all main clean debug
