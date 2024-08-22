run:
	go run ./cmd/

docker_build:
	docker build -t forum .

docker_run:
	docker run -p 8080:8080 --name forum -d forum

docker_stop:
	docker stop forum

docker_clear:
	docker rm forum