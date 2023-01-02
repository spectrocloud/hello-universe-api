
VERSION:=1.0.0

build:
	go build -o hello-universe-api

clean:
	rm -f hello-universe-api

docker-pull:
	 docker pull ghcr.io/spectrocloud/hello-universe-db:$(VERSION)

tests: build docker-pull
	docker run --detach -p 5432:5432 ghcr.io/spectrocloud/hello-universe-db:$(VERSION)
	./hello-universe-api
	newman run tests/postman_collection.json