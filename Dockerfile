FROM golang:1.12-alpine

RUN apk update \
    && apk add --no-cache git \
       openssh \
       gcc \
       cmake \
       libstdc++ libgcc g++ libquadmath musl musl-dev libgfortran \
       make \
       jpeg jpeg-dev \
       libpng libpng-dev \
       giflib giflib-dev \
       lapack \
       lapack-dev \
       openblas \
       openblas-dev \
       blas \
       ca-certificates curl wget \
       libc-dev \
    && rm -rf /var/cache/apk/*

WORKDIR /

ARG BRANCH=v19.24

RUN wget -c -q https://github.com/davisking/dlib/archive/${BRANCH}.tar.gz \
 && tar xf ${BRANCH}.tar.gz \
 && mv dlib-* dlib \
 && mkdir -p dlib/build \
 && (cd dlib/build \
    && cmake .. \
    && cmake --build . --config Release \
    && make install) \
 && rm -rf *.tar.gz /dlib/build

LABEL maintainer="Ricardo Antonio Valladares Renderos <r_a_v_r_@hotmail.com>"

WORKDIR /docker

COPY . .

COPY go.mod ./

RUN go mod download

RUN go build -o main .

EXPOSE 5000

CMD ["./main"]
