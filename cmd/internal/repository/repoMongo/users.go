package repoMongo

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"yp-diplom-2/cmd/internal/domain"
	"yp-diplom-2/cmd/internal/uuid7"

	trmmongo "github.com/avito-tech/go-transaction-manager/mongo"
	trmsql "github.com/avito-tech/go-transaction-manager/sql"
	trmcontext "github.com/avito-tech/go-transaction-manager/trm/context"
	trmman "github.com/avito-tech/go-transaction-manager/trm/manager"
)

type UserRepo struct {
	collection *mongo.Collection
	getter     *trmsql.CtxGetter
	manager    *trmman.Manager
}

func NewUserRepo(collection *mongo.Collection, getter *trmsql.CtxGetter) *UserRepo {
	return &UserRepo{
		collection: collection,
		getter:     getter,
		manager: trmman.Must(
			trmmongo.NewDefaultFactory(collection.Database().Client()),
			trmman.WithCtxManager(trmcontext.DefaultManager),
		),
	}

}

func (r *UserRepo) Create(ctx context.Context, user domain.User) error {
	err := r.manager.Do(ctx, func(ctx context.Context) error {
		_, err := r.collection.InsertOne(ctx, user)
		if IsDuplicate(err) {
			return domain.ErrUserAlreadyExists
		}
		return err
	})
	return err
}

func (r *UserRepo) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	if err := r.collection.FindOne(ctx, bson.M{"email": email, "password": password}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepo) GetByRefreshToken(ctx context.Context, refreshToken string) (domain.User, error) {
	var user domain.User
	if err := r.collection.FindOne(ctx, bson.M{
		"session.refreshToken": refreshToken,
		"session.expiresAt":    bson.M{"$gt": time.Now()},
	}).Decode(&user); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, domain.ErrUserNotFound
		}

		return domain.User{}, err
	}

	return user, nil
}

func (r *UserRepo) Verify(ctx context.Context, userID uuid7.UUID7, code string) error {
	err := r.manager.Do(ctx, func(ctx context.Context) error {
		res, err := r.collection.UpdateOne(ctx,
			bson.M{"verification.code": code, "_id": userID},
			bson.M{"$set": bson.M{"verification.verified": true, "verification.code": ""}})
		if err != nil {
			return err
		}

		if res.ModifiedCount == 0 {
			return domain.ErrVerificationCodeInvalid
		}
		return err
	})

	return err
}

func (r *UserRepo) SetSession(ctx context.Context, userID uuid7.UUID7, session domain.Session) error {
	err := r.manager.Do(ctx, func(ctx context.Context) error {
		_, err := r.collection.UpdateOne(ctx, bson.M{"_id": userID}, bson.M{"$set": bson.M{"session": session, "lastVisitAt": time.Now()}})
		return err
	})
	return err
}
