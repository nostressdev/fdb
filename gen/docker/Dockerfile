FROM golang:1.18 as builder

RUN wget --no-check-certificate https://github.com/apple/foundationdb/releases/download/6.3.22/foundationdb-clients_6.3.22-1_amd64.deb
RUN dpkg -i foundationdb-clients_6.3.22-1_amd64.deb
RUN apt update
RUN apt install -y dnsutils

WORKDIR /fdb

COPY . /fdb
