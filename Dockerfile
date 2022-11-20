FROM ubuntu:22.04

RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    curl \
    gfortran \
    git \
    gcc \
    g++ \
    wget \
    graphicsmagick \
    libgraphicsmagick1-dev \
    libatlas-base-dev \
    libavcodec-dev \
    libavformat-dev \
    libboost-all-dev \
    libgtk2.0-dev \
    libjpeg-dev \
    liblapack-dev \
    libswscale-dev \
    pkg-config \
    software-properties-common \
    zip \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

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


RUN wget -P /tmp "https://dl.google.com/go/go1.19.2.linux-amd64.tar.gz"
RUN tar -C /usr/local -xzf "/tmp/go1.19.2.linux-amd64.tar.gz"
RUN rm "/tmp/go1.19.2.linux-amd64.tar.gz"
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"


LABEL maintainer="Ricardo Antonio Valladares Renderos <r_a_v_r_@hotmail.com>"

WORKDIR /docker

COPY . .

COPY go.mod ./

RUN go mod download

RUN go build -o main .

EXPOSE 5000

CMD ["./main"]
