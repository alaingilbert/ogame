PKGS = $(shell go list ./... | grep -v /vendor/ | grep -v /bindata | grep -v /cmd/c)
VERSION  = $(shell git describe)

all:
	@echo "done"

lint:
	@golint $(PKGS)

test:
	@go test $(PKGS)

serve:
	realize start

bindata-dev:
	go-bindata -debug -pkg bindata -o cmd/scripts/bindata/bindata.go -prefix "cmd/scripts/web/public/" cmd/scripts/web/public/...

bindata-prod:
	go-bindata -pkg bindata -o cmd/scripts/bindata/bindata.go -prefix "cmd/scripts/web/public/" cmd/scripts/web/public/...

build: bindata-prod
	go build -ldflags "-s -w -X main.version=$(VERSION)" -o bot cmd/scripts/main.go

build-linux: bindata-prod
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.version=0.0.0" -o bot cmd/scripts/main.go

cover:
	@mkdir -p ./coverage
	@for pkg in $(PKGS) ; do \
		go test \
			-coverpkg=$$(go list -f '{{ join .Deps "\n" }}' $$pkg | grep '^$(PACKAGE)/' | grep -v '^$(PACKAGE)/vendor/' | tr '\n' ',')$$pkg \
			-coverprofile="./coverage/`echo $$pkg | tr "/" "-"`.cover" $$pkg ;\
	done
	@gocovmerge ./coverage/*.cover > cover.out
	@go tool cover -html=cover.out

.PHONY: bindata-dev bindata-prod build build-linux serve test lint cover
