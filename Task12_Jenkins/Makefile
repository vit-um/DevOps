ifeq '$(findstring ;,$(PATH))' ';'
    detected_OS := windows
	detected_arch := amd64
else
    detected_OS := $(shell uname | tr '[:upper:]' '[:lower:]' 2> /dev/null || echo Unknown)
    detected_OS := $(patsubst CYGWIN%,Cygwin,$(detected_OS))
    detected_OS := $(patsubst MSYS%,MSYS,$(detected_OS))
    detected_OS := $(patsubst MINGW%,MSYS,$(detected_OS))
	detected_arch := $(shell dpkg --print-architecture 2>/dev/null || amd64)
endif

#colors:
B = \033[1;94m#   BLUE
G = \033[1;92m#   GREEN
Y = \033[1;93m#   YELLOW
R = \033[1;31m#   RED
M = \033[1;95m#   MAGENTA
K = \033[K#       ERASE END OF LINE
D = \033[0m#      DEFAULT
A = \007#         BEEP

APP=$(shell basename $(shell git remote get-url origin))
REGESTRY := ghcr.io/vit-um
VERSION=$(shell git describe --tags --abbrev=0 --always)-$(shell git rev-parse --short HEAD)
TARGETARCH=amd64 
TARGETOS=${detected_OS}

format:
	gofmt -s -w ./

get:
	go get

lint:
	golint

test:
	go test -v

build: format get
	@printf "$GDetected OS/ARCH: $R$(detected_OS)/$(detected_arch)$D\n"
	CGO_ENABLED=0 GOOS=$(detected_OS) GOARCH=$(detected_arch) go build -v -o kbot -ldflags "-X="github.com/vit-um/kbot/cmd.appVersion=${VERSION}

linux: format get
	@printf "$GTarget OS/ARCH: $Rlinux/$(detected_arch)$D\n"
	CGO_ENABLED=0 GOOS=linux GOARCH=$(detected_arch) go build -v -o kbot -ldflags "-X="github.com/vit-um/kbot/cmd.appVersion=${VERSION}
	docker build --build-arg name=linux -t ${REGESTRY}/${APP}:${VERSION}-linux-$(detected_arch) .

windows: format get
	@printf "$GTarget OS/ARCH: $Rwindows/$(detected_arch)$D\n"
	CGO_ENABLED=0 GOOS=windows GOARCH=$(detected_arch) go build -v -o kbot -ldflags "-X="github.com/vit-um/kbot/cmd.appVersion=${VERSION}
	docker build --build-arg name=windows -t ${REGESTRY}/${APP}:${VERSION}-windows-$(detected_arch) .

darwin:format get
	@printf "$GTarget OS/ARCH: $Rdarwin/$(detected_arch)$D\n"
	CGO_ENABLED=0 GOOS=darwin GOARCH=$(detected_arch) go build -v -o kbot -ldflags "-X="github.com/vit-um/kbot/cmd.appVersion=${VERSION}
	docker build --build-arg name=darwin -t ${REGESTRY}/${APP}:${VERSION}-darwin-$(detected_arch) .

arm: format get
	@printf "$GTarget OS/ARCH: $R$(detected_OS)/arm$D\n"
	CGO_ENABLED=0 GOOS=$(detected_OS) GOARCH=arm go build -v -o kbot -ldflags "-X="github.com/vit-um/kbot/cmd.appVersion=${VERSION}
	docker build --build-arg name=arm -t ${REGESTRY}/${APP}:${VERSION}-$(detected_OS)-arm .

image:
	docker build . -t ${REGESTRY}/${APP}:${VERSION}-${detected_OS}-${TARGETARCH} --build-arg TARGETOS=${detected_OS} --build-arg TARGETARCH=${TARGETARCH}

push:
	docker push ${REGESTRY}/${APP}:${VERSION}-${TARGETOS}-${TARGETARCH}

dive: image
	IMG1=$$(docker images -q | head -n 1); \
	CI=true docker run -ti --rm -v /var/run/docker.sock:/var/run/docker.sock wagoodman/dive --ci --lowestEfficiency=0.99 $${IMG1}; \
	IMG2=$$(docker images -q | sed -n 2p); \
	docker rmi $${IMG1}; \
	docker rmi $${IMG2}

clean:
	@rm -rf kbot; \
	IMG1=$$(docker images -q | head -n 1); \
	if [ -n "$${IMG1}" ]; then  docker rmi -f $${IMG1}; else printf "$RImage not found$D\n"; fi
