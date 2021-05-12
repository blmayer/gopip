FROM golang:1.16 as build

ADD . /root/

RUN cd /root && go build -v

FROM python:3.9

RUN apt update && apt install zip

COPY --from=build /root/gopip /bin/gopip

CMD ["gopip"]
