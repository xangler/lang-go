FROM migrate:v1.0.0

COPY migrations /migrations
COPY migrations.sh /
 
RUN mkdir -p /work/bin/
RUN mkdir -p /work/config/
WORKDIR /work/

COPY go-leader-server.out bin/
COPY go-proxy-server.out bin/
COPY go-web-server.out bin/