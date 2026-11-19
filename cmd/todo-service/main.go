package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"example.com/devops-lab6/internal/api"
	pb "example.com/devops-lab6/proto"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("grpc.addr", ":9090", "gRPC listen address")
	flag.Parse()

	e := makeEndpoints(service{api.NewMemory[Task]()})
	router := mux.NewRouter()
	api.Route(router, "POST", "/tasks", e.create, api.JSON[Task])
	api.Route(router, "GET", "/tasks/{id}", e.getByID, api.ID)
	api.Route(router, "PATCH", "/tasks/{id}/done", e.markDone, api.ID)

	listener, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatal(err)
	}
	server := grpc.NewServer()
	pb.RegisterTodoServer(server, newGRPC(e))
	reflection.Register(server)
	go func() {
		log.Printf("todo-service gRPC %s", *grpcAddr)
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()
	log.Printf("todo-service HTTP %s", *httpAddr)
	log.Fatal(http.ListenAndServe(*httpAddr, router))
}
