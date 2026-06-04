FROM golang:1.21 AS builder

ENV GOTOOLCHAIN=auto

ARG project=parse
ARG fuzzTarget=number-fuzz

WORKDIR /src
ADD . ${project}
WORKDIR /src/${project}

RUN go test -c -fuzz=FuzzNumber -fuzztime=1s \
    -o /${fuzzTarget}.libfuzzer \
    ./tests/number

FROM ubuntu:22.04
COPY --from=builder /number-fuzz.libfuzzer /
