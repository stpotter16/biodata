shell:
	nix develop -c $$SHELL

server:
	@go run cmd/server/main.go
