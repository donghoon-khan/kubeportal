APPNAME := kube-portal
ROOTDIR := $(shell /bin/pwd)
BUILD_DIST := $(ROOTDIR)/dist
BACKEND := $(ROOTDIR)/src/app/backend
FRONTEND := $(ROOTDIR)/src/app/frontend

GOCMD := $(shell which go)
GOCLEAN :=$(GOCMD) clean
GOBUILD:=$(GOCMD) build
GOINSTALL:=$(GOCMD) install
GOTEST:=$(GOCMD) test
GOMOD:=$(GOCMD) mod

check_dist:
	mkdir -p $(BUILD_DIST)

mod:
	$(GOMOD) tidy

test:
	$(GOTEST) ./...

build/backend: check_dist
	$(GOBUILD) -o $(BUILD_DIST)/$(APPNAME) $(BACKEND)/*.go

build/frontend: check_dist
	cd $(FRONTEND); npm run build;
	mv $(FRONTEND)/build $(BUILD_DIST)/webapp

build: build/backend build/frontend

all: mod test build

clean:
	$(GOCLEAN) ./...
	rm -rf $(BUILD_DIST)
