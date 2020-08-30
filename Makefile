APP=urlshortener
APP_VERSION:=0.1
APP_COMMIT:=$(shell git rev-parse HEAD)
APP_EXECUTABLE="./out/$(APP)"
ALL_PACKAGES=$(shell go list ./... | grep -v "vendor")

setup: copy-config init-db migrate test

init-db:
	psql -c "create user shortener_user superuser password 'shortener_password';" -U postgres
	psql -c "create database shortener_db owner=shortener_user" -U postgres

compile:
	mkdir -p out/
	go build -ldflags "-X main.version=$(APP_VERSION) -X main.commit=$(APP_COMMIT)" -o $(APP_EXECUTABLE) cmd/*.go

build: deps compile

serve: build
	$(APP_EXECUTABLE) serve

tidy:
	go mod tidy

deps:
	go mod download

check: fmt vet lint

fmt:
	go fmt $(ALL_PACKAGES)

vet:
	go vet $(ALL_PACKAGES)

lint:
	golint $(ALL_PACKAGES)

clean:
	rm -rf out/

copy-config:
	cp .env.sample local.env

test:
	go clean -testcache
	go test ./...

test-cover-html:
	go clean -testcache
	mkdir -p out/
	go test ./... -coverprofile=out/coverage.out
	go tool cover -html=out/coverage.out

ci-test: copy-config init-db migrate test

docker-build:
	docker build -t nsnikhil/$(APP):$(APP_VERSION) .
	docker rmi -f $$(docker images -f "dangling=true" -q)

migrate: build
	$(APP_EXECUTABLE) migrate

rollback: build
	$(APP_EXECUTABLE) rollback