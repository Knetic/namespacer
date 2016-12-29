default: build
all: package

export GOPATH=$(CURDIR)/
export GOBIN=$(CURDIR)/.temp/

init: clean
	go get ./...

build: init
	go build -o ./.output/namespacer .

test:
	go test
	go test -bench=.

clean:
	@rm -rf ./.output/

fmt:
	@go fmt .
	@go fmt ./src/namespacer

dist: build test

	export GOOS=linux; \
	export GOARCH=amd64; \
	go build -o ./.output/namespacer64 .

	export GOOS=linux; \
	export GOARCH=386; \
	go build -o ./.output/namespacer32 .

	export GOOS=darwin; \
	export GOARCH=amd64; \
	go build -o ./.output/namespacer_osx .

	export GOOS=windows; \
	export GOARCH=amd64; \
	go build -o ./.output/namespacer.exe .

package: versionTest fpmTest dist

	fpm \
		--log error \
		-s dir \
		-t deb \
		-v $(NAMESPACER_VERSION) \
		-n namespacer \
		./.output/namespacer64=/usr/local/bin/namespacer \
		./docs/namespacer.7=/usr/share/man/man7/namespacer.7 \
		./autocomplete/namespacer=/etc/bash_completion.d/namespacer

	fpm \
		--log error \
		-s dir \
		-t deb \
		-v $(NAMESPACER_VERSION) \
		-n namespacer \
		-a i686 \
		./.output/namespacer32=/usr/local/bin/namespacer \
		./docs/namespacer.7=/usr/share/man/man7/namespacer.7 \
		./autocomplete/namespacer=/etc/bash_completion.d/namespacer

	@mv ./*.deb ./.output/

	fpm \
		--log error \
		-s dir \
		-t rpm \
		-v $(NAMESPACER_VERSION) \
		-n namespacer \
		./.output/namespacer64=/usr/local/bin/namespacer \
		./docs/namespacer.7=/usr/share/man/man7/namespacer.7 \
		./autocomplete/namespacer=/etc/bash_completion.d/namespacer
	fpm \
		--log error \
		-s dir \
		-t rpm \
		-v $(NAMESPACER_VERSION) \
		-n namespacer \
		-a i686 \
		./.output/namespacer32=/usr/local/bin/namespacer \
		./docs/namespacer.7=/usr/share/man/man7/namespacer.7 \
		./autocomplete/namespacer=/etc/bash_completion.d/namespacer

	@mv ./*.rpm ./.output/

fpmTest:
ifeq ($(shell which fpm), )
	@echo "FPM is not installed, no packages will be made."
	@echo "https://github.com/jordansissel/fpm"
	@exit 1
endif

versionTest:
ifeq ($(NAMESPACER_VERSION), )

	@echo "No 'NAMESPACER_VERSION' was specified."
	@echo "Export a 'NAMESPACER_VERSION' environment variable to perform a package"
	@exit 1
endif
