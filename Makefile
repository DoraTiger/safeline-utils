NAME=safeline-utils
REPO=DoraTiger/safeline-utils
MAIN_ENTRY=cmd/app/main.go
VERSION=$(shell git describe --always --tags)
BUILD=$(shell date +%FT%T%z)
BUILD_DIR=build
RELEASE_DIR=release
GO_BUILD=CGO_ENABLED=0 go build -mod=vendor -trimpath -ldflags '-w -s -X "github.com/DoraTiger/safeline-utils/version.Version=${VERSION}" \
		-X "github.com/DoraTiger/safeline-utils/version.Build=${BUILD}" -X "github.com/DoraTiger/safeline-utils/version.Repo=${REPO}"'

.PHONY: all clean release darwin-amd64 darwin-arm64 linux-386 linux-amd64 linux-arm linux-mips linux-mipsle linux-mips64 linux-mips64le freebsd-386 freebsd-amd64 windows-386 windows-amd64 windows-arm

PLATFORM_LIST = \
	darwin-amd64 \
	darwin-arm64 \
	linux-386 \
	linux-amd64 \
	linux-arm \
	linux-mips \
	linux-mipsle \
	linux-mips64 \
	linux-mips64le \
	freebsd-386 \
	freebsd-amd64 \
	windows-386 \
	windows-amd64 \
	windows-arm

all: clean $(PLATFORM_LIST)

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

darwin-arm64:
	GOARCH=arm64 GOOS=darwin $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

linux-386:
	GOARCH=386 GOOS=linux $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

linux-arm:
	GOARCH=arm64 GOOS=linux $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

linux-mips:
	GOARCH=mips GOOS=linux $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

linux-mipsle:
	GOARCH=mipsle GOOS=linux $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

linux-mips64:
	GOARCH=mips64 GOOS=linux $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

linux-mips64le:
	GOARCH=mips64le GOOS=linux $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

freebsd-386:
	GOARCH=386 GOOS=freebsd $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

freebsd-amd64:
	GOARCH=amd64 GOOS=freebsd $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME) ${MAIN_ENTRY}

windows-386:
	GOARCH=386 GOOS=windows $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME).exe ${MAIN_ENTRY}

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME).exe ${MAIN_ENTRY}

windows-arm:
	GOARCH=arm GOOS=windows $(GO_BUILD) -o $(BUILD_DIR)/$@/$(NAME).exe ${MAIN_ENTRY}

release: all
	bash scripts/release.sh $(NAME) $(BUILD_DIR) $(RELEASE_DIR)

clean:
	rm -rf $(BUILD_DIR)/*