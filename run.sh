docker build --tag docker-gs-ping .
docker run --publish 8080:8080 docker-gs-ping