# Build the frontend
FROM node as build_frontend

RUN mkdir -p /gocan/app/ui
WORKDIR /gocan

COPY . /gocan
RUN make frontend

#################

# Build the Go Binary.
FROM golang:1.17 as build_cli
ARG BUILD_REF

RUN mkdir /gocan
COPY go.* /gocan/
COPY --from=build_frontend /gocan/business/core/ui /gocan/business/core/ui/
WORKDIR /gocan
RUN go mod download

COPY . /gocan
RUN make backend

#################

FROM ubuntu:latest
ARG BUILD_DATE
ARG BUILD_REF
RUN mkdir /gocan
#RUN mkdir -p /root/.gocan
RUN apt update
RUN apt install -y curl git
RUN useradd -ms /bin/bash gocan
USER gocan
COPY --from=build_cli /gocan/bin/gocan /home/gocan/
WORKDIR /home/gocan
RUN mkdir /home/gocan/.gocan
RUN ./gocan setup-db
RUN ./gocan start-db
RUN ./gocan stop-db
EXPOSE 80
CMD ./gocan start-db ; ./gocan ui -p 80

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="gocan" \
      org.opencontainers.image.authors="Fouad Hamdi" \
      org.opencontainers.image.source="https://github.com/fouadh/gocan" \
      org.opencontainers.image.revision="${BUILD_REF}"