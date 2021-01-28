APPNAME := kube-portal
BUILD_DIST := dist
BACKEND := src/app/backend

.PHONY: clean docs $(BACKEND)

check_swagger:
	which swag || GO11MODULE=off go get -u github.com/swaggo/swag/cmd/swag

docs: check_swagger
	GO11MODULE=off swag init --dir $(BACKEND) --output $(BUILD_DIST)/docs --parseDependency

mod:
	go mod tidy

test:
	go test ./...

build: mod docs test
	go build -o $(BUILD_DIST)/$(APPNAME) $(BACKEND)/*.go

clean:
	rm -rf $(BUILD_DIST)