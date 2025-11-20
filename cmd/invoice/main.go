package main

import (
	"fmt"
	"net"
	"os"

	"github.com/rs/zerolog"
	"github.com/sudores/invoice-system/pkg/api/auth"
	"github.com/sudores/invoice-system/pkg/api/invoice"
	"github.com/sudores/invoice-system/pkg/api/user"
	"github.com/sudores/invoice-system/pkg/config"
	invoiceRepo "github.com/sudores/invoice-system/pkg/repo/invoice"
	userRepo "github.com/sudores/invoice-system/pkg/repo/user"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	conf, err := config.Parse()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse config")
	}
	log.Info().Msg("Starting up. Logger is ok!")

	//=========== DB Setup ===========//
	db, err := gorm.Open(postgres.Open(conf.RepoConfig.DBURL))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	log.Info().Msg("DB connection ok!")

	//=========== Invoice Setup ===========//
	invMan, err := invoiceRepo.NewInvoiceManager(db, &log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup invoice manager")
	}
	log.Info().Msg("Invoice manager created!")

	//=========== User Setup ===========//
	usrMan, err := userRepo.NewUserManager(db, &log)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup user manager")
	}
	log.Info().Msg("User manager created!")

	//=========== Jwt Setup ===========//
	jwtManager := auth.NewJwtManager(conf.Jwt)

	//=========== GRPC Setup ===========//
	lis, err := net.Listen("tcp", ":50051") // choose any free port
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen")
	}
	log.Info().Msg("Listening on port!")

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(jwtManager.UnaryInterceptor()),
	)

	//=========== GRPC Service Setup ===========//
	invSvc := invoice.NewInvoicesGrpcService(&log, invMan)
	invoice.RegisterInvoiceServiceServer(grpcServer, &invSvc)

	usrSvc := user.NewUsersGrpcService(&log, usrMan, jwtManager)
	user.RegisterUserServiceServer(grpcServer, usrSvc)

	log.Info().Msg("gRPC server listening on :50051")
	fmt.Println(grpcServer.GetServiceInfo())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
