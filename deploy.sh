#! /bin/sh

git pull origin main
docker build -t blog .
docker stop blog
docker run --name blog -p 8000:8000 -d blog
