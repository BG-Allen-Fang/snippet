FROM golang:1.15

RUN mkdir -p /usr/src/app

WORKDIR /usr/src/app
COPY . /usr/src/app

EXPOSE 4000

CMD ["go", "run", "./cmd/web"]