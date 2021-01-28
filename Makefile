APPNAME := kube-portal
ROOTDIR := $(shell /bin/pwd)
BUILD_DIST := $(ROOTDIR)/dist
BACKEND := $(ROOTDIR)/src/app/backend
FRONTEND := $(ROOTDIR)/src/app/frontend

.PHONY: clean docs all $(BACKEND)

check_dist:
	mkdir -p $(BUILD_DIST)

check_swagger:
	which swag || GO11MODULE=off go get -u github.com/swaggo/swag/cmd/swag

docs: check_swagger check_dist
	GO11MODULE=off swag init --dir $(BACKEND) --output $(BUILD_DIST)/docs --parseDependency

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

all: mod test docs build

clean:
	rm -rf $(BUILD_DIST)
