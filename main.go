package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/Cell6969/go_bank/api"
	db "github.com/Cell6969/go_bank/db/sqlc"
	"github.com/Cell6969/go_bank/gapi"
	"github.com/Cell6969/go_bank/pb"
	"github.com/Cell6969/go_bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)

	// runGinServer(config, store)
	runGrpcServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot run server HTTP:", err)
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	// Initialize api for grpc server
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// initialize grpc
	grpcServer := grpc.NewServer()

	// register protobuf into grpc
	pb.RegisterSimpleBankServer(grpcServer, server)

	// document all rpc that available
	reflection.Register(grpcServer)

	// create listener
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())

	// start grpc server
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC Server")
	}
}
