FROM golang:alpine

WORKDIR /

COPY go.mod go.sum ./

ENV GOTOOLCHAIN=auto

RUN go mod download

COPY . .

RUN go build -o Cropto-Dashboard ./main.go

COPY test.html .

EXPOSE 8080

CMD [ "./Cropto-Dashboard" ]