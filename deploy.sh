#! /bin/sh

git pull origin main
docker build -t blog .
docker stop blog
docker run --name blog -d blog
