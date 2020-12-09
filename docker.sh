#!/bin/sh
dockername="indoquran-api"

docker stop $dockername
docker rm $dockername
docker rmi $dockername
docker build -t $dockername .
docker run -it -d -p 8000:8000 --name $dockername $dockername
echo "y" | docker image prune --filter label="stage=builder"
docker logs -f $dockername
