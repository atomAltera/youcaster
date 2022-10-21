FROM golang:1.18-alpine as build

WORKDIR /opt/app

COPY . /opt/app

RUN mkdir /opt/build
RUN go build -o /opt/build/youcaster .

FROM alpine:3.13

EXPOSE 3000

RUN apk --no-cache add ca-certificates ffmpeg python3 py3-pip
RUN pip3 install --no-cache-dir --no-deps -U yt-dlp

WORKDIR /opt/app

COPY --from=build /opt/build/youcaster ./

ENTRYPOINT ["./youcaster"]