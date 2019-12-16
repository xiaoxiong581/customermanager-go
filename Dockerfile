FROM docker.io/alpine:3.10.3
MAINTAINER xiaoxiong581

EXPOSE 29080

RUN mkdir -p /home/app/log &&\
    addgroup app &&\
    adduser -D -h /home/app -G app app

ADD script/server /home/app
ADD config /home/app/config

RUN chmod -R 700 /home/app &&\
    chown -R app:app /home/app
USER app
WORKDIR /home/app

ENTRYPOINT ["./server"]