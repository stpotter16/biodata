shell:
	nix develop -c $$SHELL

server/build:
	@go build -o ./tmp/server cmd/server/main.go

server/run: server/build
	./tmp/server
