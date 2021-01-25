check_swagger:
	which swag || GO11MODULE=off go get -u github.com/swaggo/swag/cmd/swag

docs: check_swagger
	GO11MODULE=off swag init