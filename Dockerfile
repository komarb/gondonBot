FROM golang:1.16 AS builder

RUN apt-get update && apt-get install -y musl-tools

WORKDIR /build

COPY go.mod ./

RUN go mod download

COPY . .

RUN apt -y update && apt -y install libopus-dev

RUN CGO_ENABLED=1 CC=musl-gcc go build --ldflags '-linkmode external -extldflags=-static' -o ./app/gondon

FROM python:3.7-slim AS python-builder

RUN apt -y update && apt -y upgrade
RUN apt -y install wget

RUN wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp.tar.gz -O /usr/local/yt-dlp.tar.gz
RUN tar -xf /usr/local/yt-dlp.tar.gz -C /usr/local
WORKDIR /usr/local/yt-dlp
COPY --from=builder /build/pyinst.py .

RUN apt update && apt install -y apt-transport-https build-essential patchelf \
    && pip install staticx pyinstaller -r requirements.txt \
    && python3 devscripts/make_lazy_extractors.py \
    && python3 pyinst.py -F -n yt-dlp \
    && staticx ./dist/yt-dlp ./yt-dlp-static

FROM alpine:latest

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/app/gondon .
COPY --from=python-builder usr/local/yt-dlp/yt-dlp-static ./yt-dlp
COPY --from=builder /build/res /res

#COPY --from=builder /build/ffmpeg . TODO: "Static" ffmpeg does weird shit (not working). Downloading from apk
RUN apk add --no-cache --update ffmpeg

CMD ["./gondon"]