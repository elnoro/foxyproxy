VERSION 0.6
FROM golang:1.19-alpine3.15
WORKDIR /app
ENV CGO_ENABLED=0

deps:
    COPY go.mod go.sum ./
    RUN go mod download
    RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

code:
    FROM +deps
    COPY . .
    SAVE ARTIFACT /sdk

test:
    FROM +code
	RUN go test ./... -coverprofile coverage.out
	SAVE ARTIFACT coverage.out

coverage:
    COPY +test/coverage.out .
	RUN go tool cover -html=coverage.out -o coverage.html
	SAVE ARTIFACT coverage.html AS LOCAL build/coverage.html

lint:
    FROM +code
    RUN golangci-lint run -v

buildL:
    FROM +code
    ARG USERARCH
    ARG USEROS
    ENV GOARCH=$USERARCH
    ENV GOOS=$USEROS
    RUN go build -o fxpr cmd/fxpr/main.go
    SAVE ARTIFACT fxpr AS LOCAL build/fxpr

ci:
    BUILD +lint
    BUILD +test
    FROM +code
    RUN go build -o fxpr cmd/fxpr/main.go
    SAVE ARTIFACT fxpr

docker:
    ARG image=elnoro/fxpr
    COPY +ci/fxpr .
    ENTRYPOINT ["/app/fxpr"]
    SAVE IMAGE --push $image