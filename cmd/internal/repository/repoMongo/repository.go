package repoMongo

import (
	"context"
	"errors"
	"fmt"
	trmsql "github.com/avito-tech/go-transaction-manager/sql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"yp-diplom-2/cmd/internal/config"
	"yp-diplom-2/cmd/internal/repository"
)

func NewRepositories(cfg *config.Config) (*repository.Repositories, error) {
	ctx := context.Background()
	opts := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s/%s", cfg.Database.URL, cfg.Database.Port, cfg.Database.Db))
	opts.SetAuth(options.Credential{Username: cfg.Database.User, Password: cfg.Database.Password})

	client, err := mongo.Connect(ctx, opts)

	if err != nil {
		return nil, err
	}

	db := client.Database("bluff")
	collection := db.Collection(UsersCollection)

	return &repository.Repositories{
		Users: NewUserRepo(collection, trmsql.DefaultCtxGetter),
	}, nil

}

func IsDuplicate(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}
