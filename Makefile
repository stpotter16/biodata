shell:
	nix develop -c $$SHELL

server/build:
	./dev-scripts/build-server.sh

server/run: server/build
	./tmp/server

server/live:
	./dev-scripts/serve.sh

server/release:
	./dev-scripts/build-release-server.sh

server/deploy:
	./dev-scripts/deploy.sh

secrets/hmac:
	xxd -l32 /dev/urandom | xxd -r -ps | base64 | tr -d = | tr + - | tr / _


lint/go:
	./dev-scripts/lint-server.sh

lint/shell:
	./dev-scripts/check-shell.sh

lint/sql:
	./dev-scripts/check-sql.sh

lint/frontend:
	./dev-scripts/check-frontend.sh
