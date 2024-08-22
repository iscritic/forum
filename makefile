run:
 go run ./cmd/
build:
 docker build -t kokoko .
dockerRun:
 docker run -p 8080:8080 --name forum -d kokoko
sh:
 docker exec -it forum sh
stop:
 docker rm -f forum
rerun:
 make stop
 make build
 make dockerRun