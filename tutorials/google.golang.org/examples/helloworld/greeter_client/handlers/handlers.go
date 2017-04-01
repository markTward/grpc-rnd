package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	pb "github.com/markTward/gospace/tutorials/google.golang.org/examples/helloworld/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	//	addressDB   = "greeter-grpc:8000"
	addressDB   = ":8000"
	defaultName = "World!"
)

var tokens = make(chan struct{}, 100)

// HealthCheck simple
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
	log.Println(r.URL.Path, http.StatusOK)
}

// HelloWorld grpc request
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(addressDB, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(5*time.Second))

	// grpc server unreachable
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("did not connect: %v", err)
		fmt.Fprint(w, "did not connect: %v", err)
		return
	}
	defer conn.Close()

	// Greeter Client
	c := pb.NewGreeterClient(conn)

	// query string
	q := r.URL.Query()

	qname, ok := q["name"]
	var name string
	if ok {
		name = strings.Join(qname, ", ")
	} else {
		name = defaultName
	}

	rpc, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Printf("%v", err)
		fmt.Fprint(w, err)
	} else {
		log.Printf("%s?%s; grpc Message:%s", r.URL.Path, r.URL.RawQuery, rpc.Message)
		fmt.Fprint(w, rpc.Message)
	}
}