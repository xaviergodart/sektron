BIN                    := bin
BIN_NAME               := sektron
GOLANG_BIN             := go
GOLANG_111MODULE       := on
CGO_ENABLED            := 1
GOLANG_OS              := linux
GOLANG_ARCH            := amd64
GOLANG_BUILD_OPTS      := GO111MODULE=$(GOLANG_111MODULE)
GOLANG_BUILD_OPTS      += GOOS=$(GOLANG_OS)
GOLANG_BUILD_OPTS      += GOARCH=$(GOLANG_ARCH)
GOLANG_BUILD_OPTS      += CGO_ENABLED=$(CGO_ENABLED)
GOLANG_LINT            := $(BIN)/golangci-lint
ASEQDUMP_BIN           := aseqdump -p 14:0 | ts '[%H:%M:%.S]'
ASEQDUMP_NO_CLOCK_OPTS := | grep -v Clock

$(BIN):
	mkdir -p $(BIN)

$(GOLANG_LINT): $(BIN)
	GOBIN=$$(pwd)/$(BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.49.0

clean:
	-rm $(BIN)/$(BIN_NAME)
	-rm .coverage

get:
	$(GOLANG_BUILD_OPTS) $(GOLANG_BIN) get ./...

build: $(BIN) get
	$(GOLANG_BUILD_OPTS) $(GOLANG_BIN) build -o $(BIN)/sektron -tags $(BIN_NAME)
	chmod +x $(BIN)/$(BIN_NAME)

test: $(BIN) get
	$(GOLANG_BIN) test ./... -coverprofile=.coverage

checks: $(GOLANG_LINT)
	$(GOLANG_LINT) run ./...

monitor-midi:
	$(ASEQDUMP_BIN) $(ASEQDUMP_NO_CLOCK_OPTS)

monitor-midi-clock:
	$(ASEQDUMP_BIN)