OUT := RSFwd
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
VERSION := $(shell git describe --always --long)
BUILDTIME := $(shell TZ=Asia/Shanghai date -u +%Y%m%d-%H:%M)
LDFLAGS := "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILDTIME)"

all:
	go build -ldflags $(LDFLAGS) -o $(OUT) $(GO_FILES)

clean:
	rm $(OUT)
    
