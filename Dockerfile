FROM python:3.10.3-slim-bullseye

RUN apt-get -y update
RUN apt-get install -y --fix-missing \
    build-essential \
    cmake \
    gfortran \
    gcc \
    g++ \
    git \
    wget \
    curl \
    graphicsmagick \
    libgraphicsmagick1-dev \
    libatlas-base-dev \
    libavcodec-dev \
    libavformat-dev \
    libgtk2.0-dev \
    libjpeg-dev \
    liblapack-dev \
    libswscale-dev \
    pkg-config \
    python3-dev \
    python3-numpy \
    software-properties-common \
    zip \
    && apt-get clean && rm -rf /tmp/* /var/tmp/*

WORKDIR /

RUN wget -c -q "https://github.com/davisking/dlib/archive/v19.24.tar.gz"
RUN tar xf "v19.24.tar.gz" 
RUN mv dlib-* dlib
RUN mkdir -p dlib/build
RUN (cd dlib/build && cmake .. && cmake --build . --config Release && make install)
RUN rm -rf *.tar.gz /dlib/build

RUN pip3 install Flask

RUN pip3 install dlib

RUN pip3 install face_recognition

RUN pip3 install face_recognition_models

WORKDIR /docker

COPY . .

EXPOSE 5001

CMD python3 ./main.py
