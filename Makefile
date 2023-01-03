
VERSION:=1.0.0

build:
	go build -o hello-universe-api

clean:
	rm -f hello-universe-api

docker-pull-db:
	 docker pull ghcr.io/spectrocloud/hello-universe-db:$(VERSION)

docker-run-db:
	docker run --detach -p 5432:5432 --rm --name api-db ghcr.io/spectrocloud/hello-universe-db:$(VERSION)

ci-tests: build docker-pull-db docker-run-db
	sleep 3
	./hello-universe-api 2>&1 &
	newman run tests/postman_collection.json
	docker stop api-db

start-server: build docker-pull-db docker-run-db
	sleep 3
	./hello-universe-api