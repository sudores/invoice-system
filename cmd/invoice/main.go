package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/sudores/invoice-system/pkg/api"
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

	grpcServices := []api.GrpcService{invoice.NewInvoicesGrpcService(&log, invMan), user.NewUsersGrpcService(&log, usrMan, jwtManager)}

	srv := api.NewGrpcServer(grpcServices, &log, conf.Api, grpc.UnaryInterceptor(jwtManager.UnaryInterceptor()))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to serve")
	}
}
