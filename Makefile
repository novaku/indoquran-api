run-local:
	ENV ENV=local GIN_MODE=debug go run cmd/api/main.go

run-log:
	ENV ENV=local go run cmd/api/main.go >> output.log 2>&1

run-prd:
	ENV ENV=heroku go run cmd/api/main.go

test:
	go test ./... -coverprofile=coverage.out -covermode=atomic

mod:
	go mod tidy

mod-update:
	go get -u ./...
	go mod vendor

ven:
	go mod vendor

build:
	ENV ENV=heroku go build -o indoquran-api cmd/api/main.go

run:
    ./indoquran-api

dok:
	docker-compose up --build -d

dok-fresh:
	docker-compose down --volumes --rmi all
	docker-compose up --build -d

dok-drop:
	docker-compose down --volumes --rmi all