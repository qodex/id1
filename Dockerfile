FROM ubuntu:latest

RUN apt-get update && apt-get install -y tzdata
ENV TZ="Australia/Sydney"

RUN apt-get install ca-certificates -y

RUN mkdir /app
COPY ./build/api /app

RUN mkdir /api
COPY ./build/api /api
EXPOSE 8080

ENTRYPOINT [ "/app/api", "-d", "/mnt/id1db"]