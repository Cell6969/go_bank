package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/Cell6969/go_bank/api"
	db "github.com/Cell6969/go_bank/db/sqlc"
	"github.com/Cell6969/go_bank/gapi"
	"github.com/Cell6969/go_bank/pb"
	"github.com/Cell6969/go_bank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

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
	go runGatewayServer(config, store)
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
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())

	// start grpc server
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start gRPC Server")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	// Initialize api for grpc server
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	// Initiate grpcServerMux
	// Add json option for response field (the result will be snake case according to proto file)
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register Handler into grpcMux
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")
	}

	// Initiate ServerMux
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// create listener
	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start HTTP Gateway at %s", listener.Addr().String())

	// start HTTP Gateway server
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal("cannot start HTTP Gateway Server")
	}
}
