.PHONY: all
all: clean vet fmt lint test tidy gosec

.PHONY: clean
clean:
	$(call print-target)
	@go clean
	@rm -f coverage.*

.PHONY: generate
generate:
	$(call print-target)
	@go generate ./...

.PHONY: vet
vet:
	$(call print-target)
	@go vet ./...

.PHONY: fmt
fmt:
	$(call print-target)
	@go fmt ./...

.PHONY: lint
lint:
	$(call print-target)
	@golangci-lint run

.PHONY: test
test:
	$(call print-target)
	@go test -race -covermode=atomic -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

.PHONY: tidy
tidy:
	$(call print-target)
	@go mod tidy

.PHONY: gosec
gosec:
	$(call print-target)
	@gosec ./...

.PHONY: license
license:
	$(call print-target)
	@addlicense -c 'https://github.com/gio-eui' -l mit -v -s *.go

.PHONY: assets
assets:
	$(call print-target)
	# Remove all files in assets/ directory
	@rm -rf assets/*
	# Delete temp/ directory if exists and create temp/ directory
	@rm -rf temp
	@mkdir -p temp
	# Clone https://github.com/google/material-design-icons/src into temp/ directory
	@git clone -n --depth=1 --filter=tree:0 https://github.com/google/material-design-icons.git temp/
	@cd temp && git sparse-checkout set --no-cone src && git checkout
	# Copy all files from temp/src/ directory to assets/ directory
	@cp -r temp/src/* assets/
	# Remove temp/ directory
	@rm -rf temp

.PHONY: docker-build
docker-build:
	$(call print-target)
	@docker build -t svgo:latest .

.PHONY: svgo
svgo:
	$(call print-target)
	docker run -it --rm -v $(shell pwd)/assets:/shared -v $(shell pwd)/svgo.config.js:/config/svgo.config.js -v $(shell pwd)/svgo.iconvg.js:/config/svgo.iconvg.js svgo:latest svgo --config /config/svgo.config.js --folder /shared --recursive


define print-target
    @printf "Executing target: \033[36m$@\033[0m\n"
endef
