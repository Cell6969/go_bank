package main

import (
	"context"
	"database/sql"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"net"
	"net/http"

	"github.com/Cell6969/go_bank/api"
	db "github.com/Cell6969/go_bank/db/sqlc"
	_ "github.com/Cell6969/go_bank/doc/statik"
	"github.com/Cell6969/go_bank/gapi"
	"github.com/Cell6969/go_bank/pb"
	"github.com/Cell6969/go_bank/util"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config:")
	}

	if config.AppEnv == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}

	// run db migration
	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)

	// runGinServer(config, store)
	go runGatewayServer(config, store)
	runGrpcServer(config, store)
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg("cannot create migration instance:")
	}
	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("failed to migration:")
	}

	log.Info().Msg("db migration successfully")
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot run server HTTP:")
	}
}

func runGrpcServer(config util.Config, store db.Store) {
	// Initialize api for grpc server
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server")
	}
	// Add grpc log
	grpcLogerr := grpc.UnaryInterceptor(gapi.GrpcLogger)

	// initialize grpc
	grpcServer := grpc.NewServer(grpcLogerr)

	// register protobuf into grpc
	pb.RegisterSimpleBankServer(grpcServer, server)

	// document all rpc that available
	reflection.Register(grpcServer)

	// create listener
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

	// start grpc server
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Msg("cannot start gRPC Server")
	}
}

func runGatewayServer(config util.Config, store db.Store) {
	// Initialize api for grpc server
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server")
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
		log.Fatal().Msg("cannot register handler server")
	}

	// Initiate ServerMux
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	// create swagger handler
	statikFs, err := fs.New()
	if err != nil {
		log.Fatal().Msg("cannot create static fs")
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	// create listener
	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener:")
	}

	log.Info().Msgf("start HTTP Gateway at %s", listener.Addr().String())

	// start HTTP Gateway server
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatal().Msg("cannot start HTTP Gateway Server")
	}
}
