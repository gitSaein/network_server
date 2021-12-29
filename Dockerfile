FROM golang:1.16-alpine

WORKDIR /app
# Download Go modules
COPY go.mod .
COPY go.sum .
# Running go modules
RUN go mod download

COPY *.go ./
RUN go build -o /network_server
EXPOSE 8080
CMD [ "/network_server" ]