## We specify the base image we need for our go application
FROM golang:1.12.0-alpine3.9

LABEL version="1.0"
LABEL description="This project run a web server that will take your text and converted to graphical graphic representation using ASCII. Only ASCII characters are acceptable"

RUN apk add --no-cache bash

## We create an /app directory within our image that will hold our application source files
RUN mkdir /app

## We copy everything in the root directory into our /app directory
ADD . /app
## We specify that we now wish to execute any further commands inside our /app directory
WORKDIR /app

## we run go build to compile the binary executable of our Go program
RUN go build -o main .

## Our start command which kicks off, our newly created binary executable
CMD ["/app/main"]
