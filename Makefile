APPNAME := kube-portal
ROOTDIR := $(shell /bin/pwd)
BUILD_DIST := $(ROOTDIR)/dist
BACKEND := $(ROOTDIR)/src/app/backend
FRONTEND := $(ROOTDIR)/src/app/frontend

check_dist:
	mkdir -p $(BUILD_DIST)

mod:
	go mod tidy

test:
	go test ./...

build/backend: check_dist
	go build -o $(BUILD_DIST)/$(APPNAME) $(BACKEND)/*.go

build/frontend: check_dist
	cd $(FRONTEND); npm run build;
	mv $(FRONTEND)/build $(BUILD_DIST)/webapp

build: build/backend build/frontend

all: mod test build

clean:
	rm -rf $(BUILD_DIST)
