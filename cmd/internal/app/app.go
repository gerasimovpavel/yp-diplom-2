package app

import (
	"yp-diplom-2/cmd/internal/auth"
	"yp-diplom-2/cmd/internal/config"
	"yp-diplom-2/cmd/internal/hash"
	"yp-diplom-2/cmd/internal/http/handlers"
	"yp-diplom-2/cmd/internal/repository/repoPostgres"
	"yp-diplom-2/cmd/internal/server"
	"yp-diplom-2/cmd/internal/service"
)

func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		panic(err)
	}

	hasher := hash.NewSHA512Hasher(cfg.Auth.PasswordSalt)

	tokenManager, err := auth.NewManager(cfg)
	if err != nil {
		panic(err)
	}

	repos, err := repoPostgres.NewRepositories(cfg)
	if err != nil {
		panic(err)
	}

	services := service.NewServices(service.Dependencies{
		repos,
		tokenManager,
		hasher,
		cfg.Auth.JWT.AccessTokenTTL,
		cfg.Auth.JWT.RefreshTokenTTL,
	})

	handlers := handlers.NewHandler(cfg, services)

	srv := server.NewServer(cfg, handlers.Init())
	srv.Run()
}
