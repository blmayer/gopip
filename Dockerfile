FROM golang:1.16 as build

ADD main.go /root/

RUN cd /root && go build -v main.go

FROM python:3.9

COPY --from=build /root/main /bin/main

CMD ["main"]
