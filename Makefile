PKGS = $(shell go list ./... | grep -v /vendor/ | grep -v /bindata)

lint:
	@golint $(PKGS)

test:
	@go test $(PKGS)
