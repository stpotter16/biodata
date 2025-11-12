shell:
	nix develop -c $$SHELL

server/build:
	@go build -o ./tmp/server cmd/server/main.go

server/run: server/build
	./tmp/server

server/air:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o ./tmp/server cmd/server/main.go" \
	--build.delay "100" \
	--build.bin "./tmp/server" \
	--build.include_ext "go,html" \
	--build.stop_on_error "false"
