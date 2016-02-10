#!/bin/sh
docker stop postfix
docker rm postfix

docker stop postgres
docker rm  postgres

docker stop rabbitmq
docker stop redis

docker rm rabbitmq
docker rm redis

docker stop qiwi_homolog
docker rm qiwi_homolog

docker stop nginx
docker rm nginx

docker run  -e maildomain=relay.qiwi.br -e smtp_user=qiwi:vacaloca69  --name postfix -d catatnight/postfix
docker run -v  /var/lib/postgresql/data/:/var/lib/postgresql/data/  --name postgres -e POSTGRES_PASSWORD=invq1w2e3r4 -d postgres
docker run -d --name rabbitmq rabbitmq
docker run -d --name redis redis
docker run -d  --name='qiwi_homolog' --link postfix:postfix --link postgres:postgres --link redis:redis --link rabbitmq:rabbitmq -d qiwi_homolog  /root/start.sh

docker run --name nginx -p 80:80 -v /QIWI/web/password:/QIWI/web/password -v /QIWI/nginx/nginx.conf:/etc/nginx/nginx.conf:ro --link qiwi_homolog:ws -d nginx


