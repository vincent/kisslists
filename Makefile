.PHONY: build

clean:
	@rm -Rf dist

install:
	@go get ./...
	@go mod vendor
	@echo "[OK] Installed dependencies"

generate:
	@go generate ./...
	./node_modules/.bin/eslint static/frontend.js 
	@echo "[OK] Files added to embed box"

build: clean generate
	@go build -o ./dist/kisslists main.go
	@echo "[OK] App binary was created"

full: clean install generate
	CGO_ENABLED=1 go build -a -ldflags '-linkmode external -extldflags "-static"' -o ./dist/kisslists main.go
	@echo "[OK] App binary was created"

run:
	@./dist/kisslists