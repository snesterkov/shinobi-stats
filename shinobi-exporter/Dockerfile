FROM golang:1.18-alpine

RUN apk add make bash git
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN make build 
EXPOSE 9765

CMD [ "./bin/exporter" ]