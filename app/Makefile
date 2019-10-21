default: all
all: test docker

test:
	go test -covermode=count -coverprofile=count.out -v

clean:
	-rm app

build: app

app:
	CGO_ENABLED=0 go build -v -o app

docker:
	docker build -t jspc/echo-server:latest -t jspc/echo-server:$(shell date +"%s") .

publish: docker
	docker push jspc/echo-server
