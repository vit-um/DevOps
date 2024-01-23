FROM --platform=${BUILDPLATFORM} golang:1.16.6 as builder
ARG APP_BUILD_INFO=$APP_BUILD_INFO
WORKDIR /go/src/app
COPY src/ .
RUN export GOPATH=/go
RUN go get

FROM builder AS build
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o app -a -installsuffix cgo -ldflags "-X main.Version=$APP_BUILD_INFO" -v ./...

FROM golangci/golangci-lint:v1.27-alpine AS lint-base

FROM builder AS lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/root/.cache/golangci-lint \
  GO111MODULE=on golangci-lint run --disable-all -E typecheck ./...

FROM builder AS unit-test
RUN go test -v 

FROM scratch AS bin
WORKDIR /
COPY --from=build /go/src/app/app .
ENTRYPOINT ["/app"]
