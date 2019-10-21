default: all
all: docker

docker:
	docker build -t jspc/nginx:latest -t jspc/nginx:$(shell date +"%s") .

publish: docker
	docker push jspc/nginx
