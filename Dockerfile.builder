FROM docker.io/library/golang:1.23 as builder

WORKDIR /usr/src/app

RUN groupadd -g 1000 builder && useradd -u 1000 -ms /bin/bash -g builder builder
RUN chown -R builder:builder /usr/src/app

ADD ./ /usr/src/app

RUN chown -R builder:builder /usr/src/app

RUN git config --global --add safe.directory /usr/src/app

USER builder
