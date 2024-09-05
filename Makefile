.PHONY: license


VERSION:=1.1.0

build:
	go build -o hello-universe-api

clean:
	rm -f hello-universe-api

docker-pull-db:
	 docker pull ghcr.io/spectrocloud/hello-universe-db:$(VERSION)

docker-run-db:
	docker run --detach -p 5432:5432 --rm --name api-db ghcr.io/spectrocloud/hello-universe-db:$(VERSION)


stop-db:
	docker stop api-db
	docker rm api-db

ci-tests: build docker-pull-db docker-run-db
	sleep 3
	./hello-universe-api 2>&1 &
	newman run tests/postman_collection.json
	docker stop api-db

start-server: build docker-pull-db docker-run-db
	sleep 3
	./hello-universe-api


tests: docker-run-db
	sleep 3
	go test ./... -covermode=count -coverprofile=coverage.out
	go tool cover -func=coverage.out -o=coverage.out
	docker stop api-db


view-coverage: ## View the code coverage
	@echo "Viewing the code coverage"
	go tool cover -html=coverage.out

license: ## Adds a license header to all files. Reference https://github.com/hashicorp/copywrite to learn more.
	@echo "Applying license headers..."
	 copywrite headers	