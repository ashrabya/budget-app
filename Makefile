.PHONY: run build clean tidy docker-build docker-run

APP_NAME=budget-app
BINARY=./bin/$(APP_NAME)

run: tidy
	go run .

build: tidy
	mkdir -p bin
	go build -o $(BINARY) .

clean:
	rm -rf bin budget-data.json

tidy:
	go mod tidy

test:
	go test ./...

docker-build:
	docker build -t $(APP_NAME):latest .

docker-run:
	docker run -p 8080:8080 -v $(PWD)/data:/data \
		-e DATA_FILE=/data/budget-data.json \
		$(APP_NAME):latest
