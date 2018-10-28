FROM golang
WORKDIR src/
RUN go get github.com/gorilla/mux
RUN go get github.com/gocql/gocql
RUN go get github.com/scylladb/gocqlx

RUN go intall main
EXPOSE 8081